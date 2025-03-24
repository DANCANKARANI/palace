package service

import (
	"github.com/dancankarani/palace/controllers/user"
	"github.com/dancankarani/palace/model"
	"github.com/gofiber/fiber/v2"
)


func SetServicesRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/services")
	auth.Get("/all",model.GetAllServicesHandler)
	
	//protected routes
	productGroup := auth.Group("/",user.JWTMiddleware)
	productGroup.Post("/",model.CreateService)
	productGroup.Get("/",model.GetService)
	productGroup.Patch("/:id",model.UpdateServiceHandler)
	productGroup.Delete("/:id",model.DeleteServiceHandler)
}