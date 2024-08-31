package clothes

import (
	"github.com/dancankarani/palace/controllers/clothe"
	"github.com/dancankarani/palace/controllers/user"
	"github.com/gofiber/fiber/v2"
)

func SetClothesRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/clothes")
	//protected routes
	clotheGroup := auth.Group("/",user.JWTMiddleware)
	clotheGroup.Post("/",clothe.AddClotheHandler)
	clotheGroup.Patch("/:id",clothe.UpdateClotheHandler)
	clotheGroup.Delete("/:id",clothe.DeleteClotheHandler)
	clotheGroup.Get("/all",clothe.GetAllClothesHandler)
}