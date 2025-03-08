package user

import (
	"log"
	"strings"

	middleware "github.com/dancankarani/palace/middleware"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
    // Check for token in cookies first
    tokenString := c.Cookies("Authorization")
    log.Println(tokenString)
    // If not found in cookies, check the Authorization header
    if tokenString == "" {
        authHeader := c.Get("Authorization")
        if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
            tokenString = strings.TrimPrefix(authHeader, "Bearer ")
        }
    }

    // If token is still not found, return unauthorized error
    if tokenString == "" {
        log.Println("missing jwt")
        return utilities.ShowError(c, "unauthorized", fiber.StatusUnauthorized)
    }

    // Validate the token
    claims, err := middleware.ValidateToken(tokenString)
    if err != nil {
        log.Println(err.Error())
        return utilities.ShowError(c, "unauthorized", fiber.StatusUnauthorized)
    }
    //get ipd address and store in context
    ip := c.IP()
    c.Locals("ip_address", ip)
    // Store the userID in context
    c.Locals("user_id", claims.UserID)
    c.Locals("role",claims.Role)
    return c.Next()
}
