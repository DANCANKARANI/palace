package carts

import (
	"github.com/dancankarani/palace/controllers/cart"
	"github.com/dancankarani/palace/controllers/user"
	"github.com/gofiber/fiber/v2"
)

func SetCartRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/cart")
	//protected routes
	cartGroup := auth.Group("/",user.JWTMiddleware)
	cartGroup.Post("/:id",cart.AddCart)
	cartGroup.Delete("/:id/remove",cart.RemoveCartItem)
}