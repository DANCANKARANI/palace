package users

import (
	"github.com/dancankarani/palace/controllers/user"
	"github.com/gofiber/fiber/v2"
)

func SetUserRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/user")
	auth.Post("/",user.CreateUserAccount)
	auth.Post("/login",user.Login)
	auth.Get("/all",user.GetAllUsersHandler)
	//protected routes
	userGroup := auth.Group("/",user.JWTMiddleware)
	
	userGroup.Get("/",user.GetOneUserHandler)
	userGroup.Put("/",user.UpdateUserHandler)
	userGroup.Post("/forgot-password",user.ForgotPassword)
	userGroup.Post("/reset-password",user.ResetPassword)
	userGroup.Post("/logout",user.Logout)
}