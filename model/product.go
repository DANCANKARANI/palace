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

func AddCloth(c *fiber.Ctx)(*Product,error){
	cloth := Product{BaseModel: BaseModel{ID: uuid.New()}}

	//get request body
	if err := c.BodyParser(&cloth); err != nil{
		log.Println("error parsing cloth body request:",err.Error())
		return nil, errors.New("error parsing request data")
	}
	url, err := utilities.SaveFile(c,"image")
	if err != nil{
		return nil, errors.New(err.Error())
	}
	cloth.ImageURL = url
	//add to database
	if err := db.Create(&cloth).Error; err != nil{
		log.Println("error adding cloth:",err.Error())
		return nil, errors.New("failed to add cloth")
	}

	return &cloth,nil
}

/*
update cloth
@parans clothe_id
*/
func UpdateClothe(c *fiber.Ctx, clothe_id uuid.UUID)(*Product, error){
	clothe := new(Product)
	body := Product{}
	//find the clothe
	err := db.First(clothe,"id = ?",clothe_id).Error
	if err != nil{
		log.Println("error finding clothe for the update:",err.Error())
		return nil, errors.New("failed to update clothe")
	}

	//get request body
	if err := c.BodyParser(&body); err != nil{
		log.Println("failed to parse request body:",err.Error())
		return nil,errors.New("failed to parse request body")
	}
	//update clothe
	if err = db.Model(clothe).Updates(&body).Error; err != nil{
		log.Println("failed to update clothe:",err.Error())
		return nil, errors.New("failed to update clothe")
	}

	//return
	return clothe,nil
}

/*
delete clothe
@params clothe_id
*/
func DeleteClothe(c *fiber.Ctx, clothe_id uuid.UUID)error{
	clothe := new(Product)
	//get clothe
	if err := db.First(clothe, "id = ?",clothe_id).Error; err != nil{
		log.Println("error finding clothe for deleting:",err.Error())
		return errors.New("failed to delete clothe")
	}

	//delete clothe
	if err := db.Delete(clothe).Error; err != nil{
		log.Println("error deleting clothe:",err.Error())
		return errors.New("failed to delete clothe")
	}

	return nil
}

/*
gets all the clothes
*/
func GetAllClothes() (*[]Product, error) {
	var clothes []Product

	// Get all clothes
	if err := db.Find(&clothes).Error; err != nil {
		log.Println("error fetching clothes:", err.Error())
		return nil, errors.New("failed to fetch clothes")
	}

	// Shuffle clothes to return different values on each refresh
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(clothes), func(i, j int) {
		clothes[i], clothes[j] = clothes[j], clothes[i]
	})

	return &clothes, nil
}

/*
gets clothes by price
@params price
*/
func GetClothesByPrice(price float64) (*[]Product, error) {
	var clothes []Product
	// Query the database for clothes with price less than or equal to the given price
	if err := db.Where("price <= ?", price).Find(&clothes).Error; err != nil {
		log.Println("error fetching clothes by price:", err.Error())
		return nil, errors.New("failed to get clothes by price")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(clothes), func(i, j int) {
		clothes[i], clothes[j] = clothes[j], clothes[i]
	})

	return &clothes, nil
}

/*
gets clothes by gender
@params gender
*/
func GetClothesByGender(gender string) (*[]Product, error) {
	var clothes []Product
	// Query the database for clothes with the specified gender
	if err := db.Where("gender = ?", gender).Find(&clothes).Error; err != nil {
		log.Println("error fetching clothes by gender:", err.Error())
		return nil, errors.New("failed to get clothes by gender")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(clothes), func(i, j int) {
		clothes[i], clothes[j] = clothes[j], clothes[i]
	})

	return &clothes, nil
}

/*
gets clothes by category
@params category
*/
func GetClothesByCategory(category string) (*[]Product, error) {
	var clothes []Product
	// Query the database for clothes with the specified category
	if err := db.Where("category = ?", category).Find(&clothes).Error; err != nil {
		log.Println("error fetching clothes by category:", err.Error())
		return nil, errors.New("failed to get clothes by category")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(clothes), func(i, j int) {
		clothes[i], clothes[j] = clothes[j], clothes[i]
	})

	return &clothes, nil
}


/*
search clothes by various attributes
@params searchQuery
*/
func SearchClothes(searchQuery string) (*[]Product, error) {
	var clothes []Product
	// Use a case-insensitive search for the search query
	searchQuery = strings.ToLower(searchQuery)
	
	// Query the database to find clothes that match the search query in the name or description
	if err := db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%").Find(&clothes).Error; err != nil {
		log.Println("error searching clothes:", err.Error())
		return nil, errors.New("failed to search clothes")
	}

	// Shuffle the clothes slice to return them in random order
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rand.Shuffle(len(clothes), func(i, j int) {
		clothes[i], clothes[j] = clothes[j], clothes[i]
	})

	return &clothes, nil
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