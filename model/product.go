package model

import (
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddProduct(c *fiber.Ctx)(*Product,error){
	user_id,_:= GetAuthUserID(c)
	log.Println(user_id)
	product := Product{BaseModel: BaseModel{ID: uuid.New()},SellerID: user_id}
	
	//get request body
	if err := c.BodyParser(&product); err != nil{
		log.Println("error parsing product body request:",err.Error())
		return nil, errors.New("error parsing request data")
	}
	url, err := utilities.SaveFile(c,"image")
	if err != nil{
		return nil, errors.New(err.Error())
	}
	product.ImageURL = url
	//add to database
	if err := db.Create(&product).Error; err != nil{
		log.Println("error adding cloth:",err.Error())
		return nil, errors.New("failed to add cloth")
	}

	return &product,nil
}

/*
update cloth
@parans clothe_id
*/
func UpdateProduct(c *fiber.Ctx, product_id uuid.UUID)(*Product, error){
	product := new(Product)
	body := Product{}
	//find the clothe
	err := db.First(product,"id = ?",product_id).Error
	if err != nil{
		log.Println("error finding product for the update:",err.Error())
		return nil, errors.New("failed to update clothe")
	}

	//get request body
	if err := c.BodyParser(&body); err != nil{
		log.Println("failed to parse request body:",err.Error())
		return nil,errors.New("failed to parse request body")
	}
	//update clothe
	if err = db.Model(product).Updates(&body).Error; err != nil{
		log.Println("failed to update clothe:",err.Error())
		return nil, errors.New("failed to update clothe")
	}

	//return
	return product,nil
}

/*
delete clothe
@params clothe_id
*/
func DeleteProduct(c *fiber.Ctx, product_id uuid.UUID)error{
	product := new(Product)
	//get clothe
	if err := db.First(product, "id = ?",product_id).Error; err != nil{
		log.Println("error finding clothe for deleting:",err.Error())
		return errors.New("failed to delete clothe")
	}

	//delete clothe
	if err := db.Delete(product).Error; err != nil{
		log.Println("error deleting clothe:",err.Error())
		return errors.New("failed to delete clothe")
	}

	return nil
}

/*
gets all the clothes
*/
func GetAllProducts() (*[]Product, error) {
	var products []Product

	// Get all clothes
	if err := db.Find(&products).Error; err != nil {
		log.Println("error fetching clothes:", err.Error())
		return nil, errors.New("failed to fetch clothes")
	}

	// Shuffle clothes to return different values on each refresh
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(products), func(i, j int) {
		products[i], products[j] = products[j], products[i]
	})

	return &products, nil
}

/*
gets clothes by price
@params price
*/
func GetProductsByPrice(price float64) (*[]Product, error) {
	var products []Product
	// Query the database for clothes with price less than or equal to the given price
	if err := db.Where("price <= ?", price).Find(&products).Error; err != nil {
		log.Println("error fetching clothes by price:", err.Error())
		return nil, errors.New("failed to get clothes by price")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(products), func(i, j int) {
		products[i], products[j] = products[j], products[i]
	})

	return &products, nil
}

/*
gets clothes by gender
@params gender
*/
func GetProductsByGender(gender string) (*[]Product, error) {
	var product []Product
	// Query the database for clothes with the specified gender
	if err := db.Where("gender = ?", gender).Find(&product).Error; err != nil {
		log.Println("error fetching clothes by gender:", err.Error())
		return nil, errors.New("failed to get clothes by gender")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(product), func(i, j int) {
		product[i], product[j] = product[j], product[i]
	})

	return &product, nil
}

/*
gets clothes by category
@params category
*/
func GetProductsByCategory(category string) (*[]Product, error) {
	var products []Product
	// Query the database for clothes with the specified category
	if err := db.Where("category = ?", category).Find(&products).Error; err != nil {
		log.Println("error fetching clothes by category:", err.Error())
		return nil, errors.New("failed to get clothes by category")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(products), func(i, j int) {
		products[i], products[j] = products[j], products[i]
	})

	return &products, nil
}


/*
search clothes by various attributes
@params searchQuery
*/
func SearchProducts(searchQuery string) (*[]Product, error) {
	var products []Product
	// Use a case-insensitive search for the search query
	searchQuery = strings.ToLower(searchQuery)
	
	// Query the database to find clothes that match the search query in the name or description
	if err := db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%").Find(&products).Error; err != nil {
		log.Println("error searching clothes:", err.Error())
		return nil, errors.New("failed to search clothes")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(products), func(i, j int) {
		products[i], products[j] = products[j], products[i]
	})

	return &products, nil
}

/*
search and filter clothes by category and sort by price
@params category, minPrice, maxPrice, sortBy
*/
func SearchAndFilterClothes(category string, minPrice, maxPrice float64, sortBy string) (*[]Product, error) {
	var clothes []Product
	query := db.Model(&Product{})

	// Apply category filter if specified
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Apply price range filter if specified
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Apply sorting if specified
	switch strings.ToLower(sortBy) {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	default:
		query = query.Order("created_at DESC") // Default sorting
	}
	// Fetch the filtered and sorted results
	if err := query.Find(&clothes).Error; err != nil {
		log.Println("error searching and filtering clothes:", err.Error())
		return nil, errors.New("failed to search and filter clothes")
	}

	// Shuffle the clothes slice to return them in random order (optional)
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(clothes), func(i, j int) {
		clothes[i], clothes[j] = clothes[j], clothes[i]
	})

	return &clothes, nil
}

// GetSellersProduct fetches all products associated with a specific seller
func GetSellersProduct(sellerID uuid.UUID) (*[]Product, error) {
	var products []Product

	// Query the database for products with the given sellerID
	result := db.Where("seller_id = ?", sellerID).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	// If no products are found, return an empty slice
	if result.RowsAffected == 0 {
		return &products, nil
	}

	return &products, nil
}
