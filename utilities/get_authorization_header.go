package utilities

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

/*
gets the jwt token from authorization header
*/

func GetJWTToken(c *fiber.Ctx)(string,error){
	// Check for token in cookies first
    tokenString := c.Cookies("Authorization")

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
        return "",ShowError(c, "unauthorized", fiber.StatusUnauthorized)
    }
	return tokenString, nil
}