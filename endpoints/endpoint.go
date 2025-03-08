package endpoints

import (
	"github.com/dancankarani/palace/routes/carts"
	"github.com/dancankarani/palace/routes/clothes"
	"github.com/dancankarani/palace/routes/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CreateEndpoint() {
	app := fiber.New()
	
	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins, change this to specific origins in production
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", 
	}))
	users.SetUserRoutes(app)
	clothes.SetClothesRoutes(app)
	carts.SetCartRoutes(app)
	//port
	app.Listen(":8000")
}