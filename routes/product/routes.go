package product

import (
	"github.com/dancankarani/palace/controllers/product"
	"github.com/dancankarani/palace/controllers/user"
	"github.com/dancankarani/palace/model"
	"github.com/gofiber/fiber/v2"
)

func SetProductsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/products")
	auth.Get("/all",products.GetAllProductsHandler)
	auth.Get("/ratings",model.GetRatings)
	auth.Get("/price",products.GetProductsByPriceHandler)
	auth.Get("/category",products.GetProductsByCategory)
	//protected routes
	productGroup := auth.Group("/",user.JWTMiddleware)
	productGroup.Get("/",products.GetSellersProductHandler)
	productGroup.Post("/",products.AddProductHandler)
	productGroup.Post("/ratings/:id",model.CreateRatings)
	productGroup.Patch("/:id",products.UpdateProductHandler)
	productGroup.Delete("/:id",products.DeleteProductHandler)
}