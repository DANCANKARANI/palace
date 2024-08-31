package utilities

import "github.com/gofiber/fiber/v2"

func ShowSuccess(c *fiber.Ctx,msg interface{},code int,data interface{})error {
	return c.JSON(fiber.Map{
		"success":"true",
		"status_code":code,
		"message":msg,
		"data":data,
	})
}
func ShowMessage(c *fiber.Ctx,msg interface{},code int)error{
	return c.JSON(fiber.Map{
		"success":"true",
		"status_code":code,
		"message":msg,
	})
}

func ShowError(c *fiber.Ctx,errorMessage string,code int)error{
	return c.JSON(fiber.Map{
		"success":"false",
		"status_code":code,
		"error": []string{errorMessage},
	})
}