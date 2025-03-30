package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//make order functionfunc MakeOrder(userID uuid.UUID, items []OrderItem, shippingAddress, paymentMethod string) (*Order, error) {
	func MakeOrder(userID uuid.UUID, items []OrderItem, shippingAddress, paymentMethod string) (*Order, error) {
		// Validate input
		if userID == uuid.Nil {
			return nil, errors.New("user ID is required")
		}
		if len(items) == 0 {
			return nil, errors.New("at least one item is required")
		}
		if shippingAddress == "" {
			return nil, errors.New("shipping address is required")
		}
		if paymentMethod == "" {
			return nil, errors.New("payment method is required")
		}
	
		// Start transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
	
		// Create order
		order := Order{
			BaseModel:       BaseModel{ID: uuid.New()},
			OrderNumber:     generateOrderNumber(),
			UserID:          userID,
			TotalAmount:     0, // Will be calculated
			PaymentStatus:   PaymentPending,
			PaymentMethod:   paymentMethod,
			ShippingAddress: shippingAddress,
			OrderStatus:     OrderProcessing,
		}
	
		// Save order first to get ID
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order: %v", err)
		}
	
		var totalAmount float64
		// Create order items and update product stock
		for _, itemReq := range items {
			// Get product details
			var product Product
			if err := tx.First(&product, "id = ?", itemReq.ProductID).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("product not found: %v", err)
			}
	
			// Check stock
			if product.Stock < itemReq.Quantity {
				tx.Rollback()
				return nil, fmt.Errorf("not enough stock for product %s", product.Name)
			}
	
			// Create order item
			orderItem := OrderItem{
				BaseModel:   BaseModel{ID: uuid.New()},
				OrderID:     order.ID,
				ProductID:   product.ID,
				Quantity:    itemReq.Quantity,
				Price:       product.Price,
				TotalPrice:  product.Price * float64(itemReq.Quantity),
			}
	
			if err := tx.Create(&orderItem).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create order item: %v", err)
			}
	
			// Update product stock
			if err := tx.Model(&Product{}).
				Where("id = ?", product.ID).
				Update("stock", gorm.Expr("stock - ?", itemReq.Quantity)).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update product stock: %v", err)
			}
	
			totalAmount += orderItem.TotalPrice
		}
	
		// Update order with total amount
		if err := tx.Model(&Order{}).
			Where("id = ?", order.ID).
			Update("total_amount", totalAmount).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order total: %v", err)
		}
	
		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("transaction commit failed: %v", err)
		}
	
		// Reload order with items
		if err := db.Preload("Items").First(&order, order.ID).Error; err != nil {
			return nil, fmt.Errorf("failed to load created order: %v", err)
		}
	
		return &order, nil
	}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

type TimeFilter struct {
	Period string `query:"period"` // today, yesterday, week, month, or custom
	From   string `query:"from"`  // custom start date (YYYY-MM-DD)
	To     string `query:"to"`    // custom end date (YYYY-MM-DD)
}

func GetOrders(c *fiber.Ctx) error {
	// Only allow admin users to access this endpoint
	if !IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - Admin access required",
		})
	}

	// Parse query parameters
	var filter TimeFilter
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Calculate time range based on filter
	var startTime, endTime time.Time
	now := time.Now()

	switch filter.Period {
	case "today":
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = now
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		startTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
		endTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, now.Location())
	case "week":
		startTime = now.AddDate(0, 0, -7)
		endTime = now
	case "month":
		startTime = now.AddDate(0, -1, 0)
		endTime = now
	case "custom":
		// Parse custom date range
		var err error
		if filter.From != "" {
			startTime, err = time.Parse("2006-01-02", filter.From)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid 'from' date format (use YYYY-MM-DD)",
				})
			}
		} else {
			startTime = time.Date(1970, 1, 1, 0, 0, 0, 0, now.Location()) // Default to beginning of time
		}

		if filter.To != "" {
			endTime, err = time.Parse("2006-01-02", filter.To)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid 'to' date format (use YYYY-MM-DD)",
				})
			}
			// Include the entire end day
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, endTime.Location())
		} else {
			endTime = now
		}
	default:
		// If no period specified, return all orders
		startTime = time.Date(1970, 1, 1, 0, 0, 0, 0, now.Location())
		endTime = now
	}

	// Get filtered orders from database
	orders, err := GetOrdersByDateRange(startTime, endTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	// Calculate summary statistics
	totalOrders := len(orders)
	var totalRevenue float64
	var totalItems int

	for _, order := range orders {
		totalRevenue += order.TotalAmount
		for _, item := range order.Items {
			totalItems += item.Quantity
		}
	}

	// Return response with orders and summary
	return c.JSON(fiber.Map{
		"meta": fiber.Map{
			"period":       filter.Period,
			"start_date":   startTime.Format("2006-01-02"),
			"end_date":     endTime.Format("2006-01-02"),
			"total_orders": totalOrders,
			"total_items":  totalItems,
			"total_revenue": totalRevenue,
		},
		"orders": orders,
	})
}

// model/order.go
func GetOrdersByDateRange(startTime, endTime time.Time) ([]Order, error) {
	var orders []Order
	err := db.Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Preload("Items").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func IsAdmin(c *fiber.Ctx) bool {
	// Implement your admin check logic here
	// Example: return c.Locals("user_role") == "admin"
	return true
}