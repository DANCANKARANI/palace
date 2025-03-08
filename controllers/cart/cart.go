package cart

import (
	"github.com/dancankarani/palace/model"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddCart(c *fiber.Ctx)error{
	product_id,_:=uuid.Parse(c.Params("id"))
	res,err := model.AddCart(c,product_id)
	if err != nil{
		return utilities.ShowError(c, err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowMessage(c,res,fiber.StatusOK)
}

func RemoveCartItem(c *fiber.Ctx)error{
	cart_item_id,_:= uuid.Parse(c.Params("id"))
	err := model.RemoveCartItem(cart_item_id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	return utilities.ShowMessage(c,"item removed successfully",fiber.StatusOK)
}

