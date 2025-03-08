package model

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddCart(c *fiber.Ctx,product_id uuid.UUID) (*CartItem, error) {
    userID, _ := GetAuthUserID(c)
    log.Println(userID)
    // Parse the request body to get CartItem data
    var cartItem CartItem
    if err := c.BodyParser(&cartItem); err != nil {
        log.Println("Error parsing cart items request:", err.Error())
        return nil, errors.New("failed to read request data")
    }

    // Find the user's cart or create a new cart if it doesn't exist
    var cart Cart
    if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Create a new cart if not found
            cart = Cart{
                BaseModel:   BaseModel{ID: uuid.New()},
                UserID:      userID,
                TotalAmount: 0, // Initialize to 0
            }
            if err := db.Create(&cart).Error; err != nil {
                log.Println("Error creating new cart:", err.Error())
                return nil, errors.New("failed to create cart")
            }
        } else {
            // Return other database errors
            log.Println("Error retrieving cart:", err.Error())
            return nil, errors.New("failed to retrieve cart")
        }
    }

    // Set the CartID for the CartItem
    cartItem.ProductID = product_id

    // Calculate TotalPrice for the cart item
    cartItem.TotalPrice = float64(cartItem.Quantity) * cartItem.Price

    // Add the cart item to the database
    cartItem = CartItem{ BaseModel:   BaseModel{ID: uuid.New()},CartID:cart.ID,ProductID: product_id,Quantity: cartItem.Quantity,Price: cartItem.Price,TotalPrice: cartItem.TotalPrice}
    if err := db.Create(&cartItem).Error; err != nil {
        log.Println("Error adding cart item:", err.Error())
        return nil, errors.New("failed to add cart item")
    }

    // Update cart total amount
    var cartItems []CartItem
    if err := db.Model(&cart).Association("Items").Find(&cartItems); err != nil {
        log.Println("Error retrieving cart items:", err.Error())
        return nil, errors.New("failed to retrieve cart items")
    }

    totalAmount := 0.0
    for _, item := range cartItems {
        totalAmount += item.TotalPrice
    }
    cart.TotalAmount = totalAmount

    if err := db.Save(&cart).Error; err != nil {
        log.Println("Error updating cart total amount:", err.Error())
        return nil, errors.New("failed to update cart total amount")
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