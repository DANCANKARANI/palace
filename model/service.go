package model

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



func CreateService(c *fiber.Ctx) error {
	// Get the authenticated user ID
	userID, err := GetAuthUserID(c)
	if err != nil {
		log.Println("error getting authenticated user ID:", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	log.Println("Authenticated user ID:", userID)

	// Parse the request body into the Service struct
	var service Service
	if err := c.BodyParser(&service); err != nil {
		log.Println("error parsing service body request:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request data",
		})
	}

	// Set the SellerID to the authenticated user's ID
	service.SellerID = userID

	// Generate a unique ID for the service
	service.ID = uuid.New()




	// Add the service to the database
	if err := db.Create(&service).Error; err != nil {
		log.Println("error adding service to database:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create service",
		})
	}

	// Return the created service with a 201 status code
	return c.Status(fiber.StatusCreated).JSON(service)
}
// GetService retrieves all services for the authenticated seller
func GetService(c *fiber.Ctx) error {
	// Get the authenticated user ID
	userID, err := GetAuthUserID(c)
	if err != nil {
		log.Println("error getting authenticated user ID:", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	log.Println("Authenticated user ID:", userID)

	// Query the database for services associated with the seller
	var services []Service
	if err := db.Where("seller_id = ?", userID).Find(&services).Error; err != nil {
		log.Println("error fetching services:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch services",
		})
	}

	// Return the list of services
	return c.Status(fiber.StatusOK).JSON(services)
}

// GetAllServicesHandler retrieves all services from the database
func GetAllServicesHandler(c *fiber.Ctx) error {
	// Extract query parameters for filtering
	category := c.Query("category") // Get the category from query parameters

	// Query the database for all services and preload the User details
	var services []Service
	query := db.Preload("User") // Preload the User details

	// Apply category filter if provided
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Execute the query
	if err := query.Find(&services).Error; err != nil {
		log.Println("error fetching services:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch services",
		})
	}

	// Return the list of services with User details
	return c.Status(fiber.StatusOK).JSON(services)
}
// UpdateServiceHandler updates an existing service
func UpdateServiceHandler(c *fiber.Ctx) error {
	// Get the authenticated user ID
	userID, err := GetAuthUserID(c)
	if err != nil {
		log.Println("error getting authenticated user ID:", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get the service ID from the request parameters
	serviceID := c.Params("id")
	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service ID is required",
		})
	}

	// Parse the request body into a map to allow partial updates
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("error parsing request body:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request data",
		})
	}

	// Find the service by ID and ensure it belongs to the authenticated seller
	var service Service
	if err := db.Where("id = ? AND seller_id = ?", serviceID, userID).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "service not found or unauthorized",
			})
		}
		log.Println("error fetching service:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch service",
		})
	}

	// Update the service with the new data
	if err := db.Model(&service).Updates(updateData).Error; err != nil {
		log.Println("error updating service:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update service",
		})
	}

	// Return the updated service
	return c.Status(fiber.StatusOK).JSON(service)
}


// DeleteServiceHandler deletes an existing service
func DeleteServiceHandler(c *fiber.Ctx) error {
	// Get the authenticated user ID
	userID, err := GetAuthUserID(c)
	if err != nil {
		log.Println("error getting authenticated user ID:", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get the service ID from the request parameters
	serviceID := c.Params("id")
	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service ID is required",
		})
	}

	// Find the service by ID and ensure it belongs to the authenticated seller
	var service Service
	if err := db.Where("id = ? AND seller_id = ?", serviceID, userID).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "service not found or unauthorized",
			})
		}
		log.Println("error fetching service:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch service",
		})
	}

	// Delete the service
	if err := db.Delete(&service).Error; err != nil {
		log.Println("error deleting service:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete service",
		})
	}

	// Return a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "service deleted successfully",
	})
}