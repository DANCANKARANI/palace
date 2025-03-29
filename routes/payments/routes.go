package payments

import (
	"github.com/dancankarani/palace/controllers/payment"
	"github.com/gofiber/fiber/v2"
)

func SetPaymentsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/")
	auth.Post("/payments", payment.InitiateSTKPush)
	auth.Post("/callback",payment.HandleCallback)
}