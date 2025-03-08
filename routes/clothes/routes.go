package clothes

import (
	"github.com/dancankarani/palace/controllers/clothe"
	"github.com/dancankarani/palace/controllers/user"
	"github.com/gofiber/fiber/v2"
)

func SetClothesRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/clothes")
	auth.Get("/all",clothe.GetAllClothesHandler)
	auth.Get("/price",clothe.GetClothesByPriceHandler)
	auth.Get("/category",clothe.GetClothesByCategory)
	//protected routes
	clotheGroup := auth.Group("/",user.JWTMiddleware)
	clotheGroup.Post("/",clothe.AddClotheHandler)
	clotheGroup.Patch("/:id",clothe.UpdateClotheHandler)
	clotheGroup.Delete("/:id",clothe.DeleteClotheHandler)
}