package model

import (
	"errors"
	"fmt"
	"time"

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

	// Calculate total amount and validate product stock
	var totalAmount float64
	for i, item := range items {
		var product Product
		if err := db.First(&product, "id = ?", item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("product not found: %v", err)
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product %s", product.Name)
		}
		items[i].Price = product.Price
		items[i].TotalPrice = product.Price * float64(item.Quantity)
		totalAmount += items[i].TotalPrice
	}

	// Create order
	order := Order{
		OrderNumber:   generateOrderNumber(),
		UserID:        userID,
		TotalAmount:   totalAmount,
		PaymentStatus: PaymentPending,
		PaymentMethod: paymentMethod,
		ShippingAddress: shippingAddress,
		OrderStatus:   OrderProcessing,
		Items:         items,
	}

	// Save order and items in a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Update product stock
		for _, item := range items {
			if err := tx.Model(&Product{}).Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	return &order, nil
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

