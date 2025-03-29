package model

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateRatings(c *fiber.Ctx) error {
    // Initialize rating and get IDs
    rating := new(Rating)
    sellerID := c.Params("id")
    userID, err := GetAuthUserID(c)
    
    // Error handling
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unauthorized",
            "message": "Could not authenticate user",
        })
    }

    // Parse request body
    if err := c.BodyParser(rating); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Bad Request",
            "message": "Could not parse rating data",
            "details": err.Error(),
        })
    }

    // Validate input
    if rating.Stars < 1 || rating.Stars > 5 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Bad Request",
            "message": "Rating must be between 1 and 5 stars",
        })
    }

    // Set required fields
    rating.SellerID = uuid.MustParse(sellerID)
    rating.UserID = userID
    rating.ID = uuid.New()

    // Save to database
    if err := db.Create(&rating).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database Error",
            "message": "Could not create rating",
            "details": err.Error(),
        })
    }

    // Return success response
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "data": rating,
    })
}

//Get the ratings
// GetRatings retrieves ratings with optional filtering
// @Summary Get ratings
// @Description Get ratings with optional filtering by seller or user
// @Tags Ratings
// @Accept json
// @Produce json
// @Param seller_id query string false "Filter by seller ID"
// @Param user_id query string false "Filter by user ID"
// @Param limit query int false "Limit results" default(10)
// @Param page query int false "Page number" default(1)
// @Success 200 {array} Rating
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ratings [get]
func GetRatings(c *fiber.Ctx) error {
    // Initialize query
    query := db.Model(&Rating{}).Preload("User") // Assuming you want to include user details

    // Get query parameters
    sellerID := c.Query("seller_id")
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    page, _ := strconv.Atoi(c.Query("page", "1"))

    // Apply filters
    if sellerID != "" {
        query = query.Where("seller_id = ?", sellerID)
    }

    // Pagination
    offset := (page - 1) * limit
    query = query.Limit(limit).Offset(offset).Order("created_at DESC")

    // Execute query
    var ratings []Rating
    if err := query.Find(&ratings).Error; err != nil {
		fmt.Println(err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			
            "error": "Database Error",
            "message": "Could not retrieve ratings",
        })
    }

    // Count total records (for pagination metadata)
    var total int64
    countQuery := db.Model(&Rating{})
    if sellerID != "" {
        countQuery = countQuery.Where("seller_id = ?", sellerID)
    }
    countQuery.Count(&total)

    // Return response with pagination metadata
    return c.JSON(fiber.Map{
        "data": ratings,
        "meta": fiber.Map{
            "total":     total,
            "page":      page,
            "limit":     limit,
            "totalPages": int(math.Ceil(float64(total) / float64(limit))),
        },
    })
}