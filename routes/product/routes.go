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
	productGroup := auth.Group("/",user.JWTMiddleware)
	productGroup.Get("/:id",products.GetSellersProductHandler)
	productGroup.Post("/",products.AddProductHandler)
	productGroup.Patch("/:id",products.UpdateProductHandler)
	productGroup.Delete("/:id",products.DeleteProductHandler)
}