package products

import (
	"strconv"
	"github.com/dancankarani/palace/model"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddProductHandler(c *fiber.Ctx)error{
	clothe, err := model.AddProduct(c)
	if err != nil{
		return utilities.ShowError(c, err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c, "successfull added the clothe", fiber.StatusOK,clothe)
}

func UpdateProductHandler(c *fiber.Ctx)error{
	id, _:=uuid.Parse(c.Params("id"))
	clothe, err := model.UpdateProduct(c,id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"clothes updated successfully",fiber.StatusOK, clothe)
}

func GetAllProductsHandler(c *fiber.Ctx)error{
	response,err := model.GetAllProducts()
	if err != nil{
		return utilities.ShowError(c, err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"clothes retrieved successfully",fiber.StatusOK,response)
}

func DeleteProductHandler(c *fiber.Ctx)error{
	id, _:=uuid.Parse(c.Params("id"))
	err := model.DeleteProduct(c,id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowMessage(c,"clothe deleted successfully",fiber.StatusOK)
}

func GetProductsByPriceHandler(c *fiber.Ctx)error{
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice", "0"), 64)
	response, err := model.GetProductsByPrice(maxPrice)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"retrieved clothes by price",fiber.StatusOK,response)
}

func GetProductsByCategory(c *fiber.Ctx)error{
	category := c.Query("categories","unisex")
	response, err := model.GetProductsByCategory(category)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"successfull retrieved clothes by category",fiber.StatusOK,response)
}

func GetSellersProductHandler(c *fiber.Ctx)error{
	id, _:=uuid.Parse(c.Params("id"))
	response,err := model.GetSellersProduct(id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowSuccess(c,"seller's products retrieved successfully",fiber.StatusOK,response)
}