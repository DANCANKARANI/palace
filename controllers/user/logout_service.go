package user

import (
	"time"
	"github.com/dancankarani/palace/moddleware"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
)

func LogoutService(c *fiber.Ctx, user_type string) error {

	//get token string
	tokenString, err := utilities.GetJWTToken(c)
	if err != nil {
		return utilities.ShowError(c, err.Error(), fiber.StatusUnauthorized)
	}

	//invalidate token
	err = middleware.InvalidateToken(tokenString)
	if err != nil {
		return utilities.ShowError(c, "failed to invalidate the token", fiber.StatusInternalServerError)
	}


	
	//set token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	//response
	return nil
}