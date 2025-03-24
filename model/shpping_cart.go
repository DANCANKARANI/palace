package model

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddCart(c *fiber.Ctx, productID uuid.UUID) (*CartItem, error) {
    userID, err := GetAuthUserID(c)
    if err != nil {
        log.Println("Error retrieving user ID:", err.Error())
        return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
    }

    // Parse the request body to get CartItem data
    var cartItem CartItem
    if err := c.BodyParser(&cartItem); err != nil {
        log.Println("Error parsing cart item request:", err.Error())
        return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request data")
    }

    // Validate cart item fields
    if cartItem.Quantity <= 0 || cartItem.Price <= 0 {
        return nil, fiber.NewError(fiber.StatusBadRequest, "quantity and price must be greater than 0")
    }

    // Start a database transaction
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Find the user's cart or create a new cart if it doesn't exist
    var cart Cart
    if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Create a new cart if not found
            cart = Cart{
                BaseModel:   BaseModel{ID: uuid.New()},
                UserID:      userID,
                TotalAmount: 0, // Initialize to 0
            }
            if err := tx.Create(&cart).Error; err != nil {
                log.Println("Error creating new cart:", err.Error())
                return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create cart")
            }
        } else {
            // Return other database errors
            log.Println("Error retrieving cart:", err.Error())
            return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve cart")
        }
    }

    // Create the cart item
    cartItem = CartItem{
        BaseModel:  BaseModel{ID: uuid.New()},
        CartID:     cart.ID,
        ProductID:  productID,
        Quantity:   cartItem.Quantity,
        Price:      cartItem.Price,
        TotalPrice: float64(cartItem.Quantity) * cartItem.Price,
    }

    // Add the cart item to the database
    if err := tx.Create(&cartItem).Error; err != nil {
        log.Println("Error adding cart item:", err.Error())
        tx.Rollback()
        return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add cart item")
    }

    // Update cart total amount
    cart.TotalAmount += cartItem.TotalPrice
    if err := tx.Save(&cart).Error; err != nil {
        log.Println("Error updating cart total amount:", err.Error())
        tx.Rollback()
        return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update cart total amount")
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        log.Println("Error committing transaction:", err.Error())
        return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to complete transaction")
    }

    return &cartItem, nil
}

/*
removes cart items
@params cart_item_id
*/
func RemoveCartItem(cart_item_id uuid.UUID)error{
	cartItem := new(CartItem)
	if err := db.First(cartItem,"id = ?",cart_item_id).Error; err != nil{
		log.Println("error getting cart item:",err.Error())
		return errors.New("failed to remove cart item")
	}

	//delete
	if err :=db.Delete(cartItem).Error; err != nil{
		log.Println("error deleting cart item:",err.Error())
		return errors.New("error removing cart item")
	}
	return nil
}

/*
update cart items
@params cart_item_id
*/
func UpdateCart(c *fiber.Ctx, cart_item_id uuid.UUID)(*CartItem,error){
	cartItem := new(CartItem)
	if err := c.BodyParser(cartItem); err != nil{
		log.Println("error parsing cart items for updates:",err.Error())
		return nil,errors.New("error updating cart")
	}
	//find cart item
	if err := db.First(cartItem,"id = ?",cart_item_id).Error; err != nil{
		log.Println("error finding cart item for update:",err.Error())
		return nil, errors.New("failed to update cart")
	}

	//update
	if err := db.Updates(cartItem).Error; err != nil{
		log.Println("error updating cart:",err.Error())
		return nil, errors.New("error updating the cart")
	}
	return cartItem,nil
}