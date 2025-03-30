package order

import (


	"github.com/dancankarani/palace/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderRequest struct {
	Items          []OrderItemRequest `json:"items"`
	ShippingAddress string            `json:"shipping_address"`
	PaymentMethod   string            `json:"payment_method"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

func MakeOrderHandler(c *fiber.Ctx) error {
	// Get the authenticated user's ID
	userID, err := model.GetAuthUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse the request body
	var req OrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Convert OrderItemRequest to OrderItem
	var items []model.OrderItem
	for _, item := range req.Items {
		items = append(items, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Call the MakeOrder function
	order, err := model.MakeOrder(userID, items, req.ShippingAddress, req.PaymentMethod)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the created order
	return c.Status(fiber.StatusCreated).JSON(order)
}

