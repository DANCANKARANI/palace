package product

import (
	"github.com/dancankarani/palace/controllers/product"
	"github.com/dancankarani/palace/controllers/user"
	"github.com/gofiber/fiber/v2"
)

func SetProductsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/products")
	auth.Get("/all",products.GetAllProductsHandler)
	auth.Get("/price",products.GetProductsByPriceHandler)
	auth.Get("/category",products.GetProductsByCategory)
	//protected routes
	clotheGroup := auth.Group("/",user.JWTMiddleware)
	clotheGroup.Post("/",products.AddProductHandler)
	clotheGroup.Patch("/:id",products.UpdateProductHandler)
	clotheGroup.Delete("/:id",products.DeleteProductHandler)
}