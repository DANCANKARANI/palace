package orders

import (
	"github.com/dancankarani/palace/controllers/order"
	"github.com/dancankarani/palace/controllers/user"
	"github.com/gofiber/fiber/v2"
)

func SetOrdersRoutes(app *fiber.App){
	auth := app.Group("/api/v1/orders")
	productGroup := auth.Group("/",user.JWTMiddleware)
	productGroup.Post("/",order.MakeOrderHandler)
}