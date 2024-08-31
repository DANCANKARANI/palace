package clothe

import (
	"strconv"

	"github.com/dancankarani/palace/model"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddClotheHandler(c *fiber.Ctx)error{
	clothe, err := model.AddCloth(c)
	if err != nil{
		return utilities.ShowError(c, err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c, "successfull added the clothe", fiber.StatusOK,clothe)
}

func UpdateClotheHandler(c *fiber.Ctx)error{
	id, _:=uuid.Parse(c.Params("id"))
	clothe, err := model.UpdateClothe(c,id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"clothes updated successfully",fiber.StatusOK, clothe)
}

func GetAllClothesHandler(c *fiber.Ctx)error{
	response,err := model.GetAllClothes()
	if err != nil{
		return utilities.ShowError(c, err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"clothes retrieved successfully",fiber.StatusOK,response)
}

func DeleteClotheHandler(c *fiber.Ctx)error{
	id, _:=uuid.Parse(c.Params("id"))
	err := model.DeleteClothe(c,id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowMessage(c,"clothe deleted successfully",fiber.StatusOK)
}

func GetClothesByPriceHandler(c *fiber.Ctx)error{
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice", "0"), 64)
	response, err := model.GetClothesByPrice(maxPrice)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"retrieved clothes by price",fiber.StatusOK,response)
}

func GetClothesByCategory(c *fiber.Ctx)error{
	category := c.Query("categories","unisex")
	response, err := model.GetClothesByCategory(category)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"successfull retrieved clothes by category",fiber.StatusOK,response)
}
