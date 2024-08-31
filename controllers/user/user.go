package user

import (
	"github.com/dancankarani/palace/model"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
)

//get one user handler
func GetOneUserHandler(c *fiber.Ctx) error {
	user,err := model.GetOneUSer(c)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"user retrieved successfully",fiber.StatusOK,user)
}

//get all users handler
func GetAllUsersHandler(c *fiber.Ctx)error{
	response,err := model.GetAllUsersDetails(c)
	if err != nil {
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError) 
	}
	return utilities.ShowSuccess(c,"users retrieved successfully",fiber.StatusOK,response)
}

//update user details handler
func UpdateUserHandler(c *fiber.Ctx)error{
	response,err := model.UpdateUser(c)
	if err != nil {
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"user updated successfully",fiber.StatusOK,response)
}

