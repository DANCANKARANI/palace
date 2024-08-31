package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
type User struct {
    BaseModel
    FullName   string `json:"full_name" gorm:"size:255"`
    Email      string `json:"email" gorm:"size:100;unique"`
    Password   string `json:"password" gorm:"size:255"`
    Address    string `json:"address" gorm:"size:255"`
    City       string `json:"city" gorm:"size:100"`
    PostalCode string `json:"postal_code" gorm:"size:20"`
    Location   string `json:"location" gorm:"size:100"`
    PhoneNumber      string `json:"phone_number" gorm:"size:20"`
    UserRole   string `json:"user_role" gorm:"size:50;default:'customer'"` // Can be 'customer', 'admin', etc.
    IsActive   bool   `json:"is_active" gorm:"default:true"`
	ResetCode  string `json:"reset_code" gorm:"size:10"`
	CodeExpirationTime time.Time	`json:"code_expiration_time"`
}


type Clothe struct {
    BaseModel
    Name        string  `json:"name" gorm:"size:255"`
    Description string  `json:"description" gorm:"type:text"`
    Price       float64 `json:"price" gorm:"type:decimal(10,2)"`
    Size        string  `json:"size" gorm:"size:50"`    // Example: S, M, L, XL
    Color       string  `json:"color" gorm:"size:50"`
    Gender      string  `json:"gender" gorm:"size:20"`  // Example: Male, Female, Unisex
    Category    string  `json:"category" gorm:"size:100"` // Example: Shirts, Trousers, etc.
    Stock       int     `json:"stock"`                  // Available stock count
    ImageURL    string  `json:"image_url" gorm:"size:255"` // URL for the clothing item image
    Brand       string  `json:"brand" gorm:"size:100"`  // Clothing brand
    Material    string  `json:"material" gorm:"size:100"` // Example: Cotton, Polyester
    IsActive    bool    `json:"is_active" gorm:"default:true"` // Active status for display
}