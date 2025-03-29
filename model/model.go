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
    FirstName        string     `json:"first_name" gorm:"size:255;not null"`
    LastName         string     `json:"last_name" gorm:"size:255;not null"`
    Email            string     `json:"email" gorm:"size:100;unique;not null" validate:"required,email"`
    Password         string     `json:"-" gorm:"size:255;not null"` // Hidden from JSON
    Address          string     `json:"address" gorm:"size:255"`
    City             string     `json:"city" gorm:"size:100"`
    PostalCode       string     `json:"postal_code" gorm:"size:20"`
    Location         string     `json:"location" gorm:"size:100"`
    PhoneNumber      string     `json:"phone_number" gorm:"size:20" validate:"omitempty,numeric"`
    UserRole         string     `json:"user_role" gorm:"size:50;default:'customer';not null" validate:"oneof=customer admin seller"`
    IsActive         bool       `json:"is_active" gorm:"default:true"`
    ResetCode        string     `json:"-" gorm:"size:10"` // Hidden from JSON
    CodeExpirationTime time.Time `json:"-"` // Hidden from JSON
    
    // Relationships
    Services         []Service  `json:"services,omitempty" gorm:"foreignKey:SellerID"`
    Orders           []Order    `json:"orders,omitempty" gorm:"foreignKey:UserID"`
    Products         []Product  `json:"products,omitempty" gorm:"foreignKey:SellerID"`
    
    // As a seller receiving ratings
    SellerRatings    []Rating   `json:"seller_ratings,omitempty" gorm:"foreignKey:SellerID"`
    
    // As a buyer giving ratings
    GivenRatings     []Rating   `json:"given_ratings,omitempty" gorm:"foreignKey:UserID"`
}


type Product struct {
    BaseModel
    Name        string    `json:"name" gorm:"size:255"`
    Description string    `json:"description" gorm:"type:text"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
    Category    string    `json:"category" gorm:"size:100"` // Example: Shirts, Trousers, etc.
    Stock       int       `json:"stock"`                  // Available stock count
    ImageURL    string    `json:"image_url" gorm:"size:255"` // URL for the clothing item image
    IsActive    bool      `json:"is_active" gorm:"default:true"` // Active status for display
    SellerID    uuid.UUID `json:"seller_id" gorm:"type:uuid;index"` // Foreign key to associate with Seller
    User      User      `gorm:"foreignKey:SellerID;references:ID"` // Relationship to User
    OrderItems  []OrderItem `gorm:"foreignKey:ProductID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
    CartItems   []CartItem  `gorm:"foreignKey:ProductID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
}

type Service struct {
    BaseModel
    Name        string    `json:"name" gorm:"size:255"`
    Description string    `json:"description" gorm:"type:text"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
    Category    string    `json:"category" gorm:"size:100"` // Example: Graphic Design, Tutoring, etc.
    IsActive    bool      `json:"is_active" gorm:"default:true"` // Active status for display
    SellerID    uuid.UUID `json:"seller_id" gorm:"type:uuid;index"` // Foreign key to associate with Seller
    User        User      `gorm:"foreignKey:SellerID;references:ID"` // Relationship to User
}

type Order struct {
	BaseModel
	OrderNumber   string          `json:"order_number" gorm:"size:100;unique;"`
	UserID        uuid.UUID       `json:"user_id" gorm:"index;"`
	User          User            `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
	TotalAmount   float64         `json:"total_amount" gorm:"type:decimal(10,2)"`
	PaymentStatus PaymentStatus   `json:"payment_status" gorm:"size:50"`
	PaymentMethod string          `json:"payment_method" gorm:"size:50"`
	ShippingAddress string        `json:"shipping_address" gorm:"type:text"`
	OrderStatus   OrderStatus     `json:"order_status" gorm:"size:50"`
	DeliveredAt   *time.Time     `json:"delivered_at"`
	Items         []OrderItem     `json:"items" gorm:"foreignKey:OrderID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
}

type OrderItem struct {
	BaseModel
	OrderID     uuid.UUID  `json:"order_id" gorm:"index;"`
	Order       Order      `json:"order" gorm:"foreignKey:OrderID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	ProductID   uuid.UUID  `json:"product_id" gorm:"index;"`
	Product     Product    `json:"product" gorm:"foreignKey:ProductID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
	Quantity    int        `json:"quantity" gorm:"int"`
	Price       float64    `json:"price" gorm:"type:decimal(10,2)"`
	TotalPrice  float64    `json:"total_price" gorm:"type:decimal(10,2)"` // Optional: Can be calculated dynamically
}

// PaymentStatus and OrderStatus enums
type PaymentStatus string

const (
	PaymentPending PaymentStatus = "Pending"
	PaymentPaid    PaymentStatus = "Paid"
	PaymentFailed  PaymentStatus = "Failed"
)

type OrderStatus string

const (
	OrderProcessing OrderStatus = "Processing"
	OrderShipped    OrderStatus = "Shipped"
	OrderDelivered  OrderStatus = "Delivered"
	OrderCancelled  OrderStatus = "Cancelled"
)
type Cart struct {
	BaseModel
	UserID        uuid.UUID     `json:"user_id" gorm:"type:varchar(36);"` // Reference to the user who owns the cart
	User          User          `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	TotalAmount   float64       `json:"total_amount" gorm:"type:decimal(10,2)"`
	Items         []CartItem    `json:"items" gorm:"foreignKey:CartID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
}

type CartItem struct {
	BaseModel
	CartID      uuid.UUID  `json:"cart_id" gorm:"type:varchar(36);"` // Foreign key to Cart
	Cart        Cart       `json:"cart" gorm:"foreignKey:CartID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	ProductID   uuid.UUID  `json:"product_id" gorm:"type:varchar(36);"` // Foreign key to Product
	Product     Product    `json:"product" gorm:"foreignKey:ProductID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
	Quantity    int        `json:"quantity" gorm:"int"`
	Price       float64    `json:"price" gorm:"type:decimal(10,2);"`
	TotalPrice  float64    `json:"total_price" gorm:"type:decimal(10,2);"`
}

type Payment struct {
	ID              uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey;"` // Unique identifier for the payment
	BillingID       uuid.UUID `json:"billing_id" gorm:"type:varchar(36)"` // Foreign key to the Billing table
	CustomerID       uuid.UUID `json:"patient_id" gorm:"type:varchar(36)"` // Foreign key to the Patients table
	Cost            float64   `json:"cost" gorm:"type:decimal(10,2);"`        // Cost of the payment
	PaymentMethod   string    `json:"payment_method" gorm:"type:varchar(50);"` // Payment method (e.g., M-Pesa, Credit Card)
	TransactionID   string    `json:"transaction_id" gorm:"type:varchar(100);"` // Transaction ID from the payment gateway
	PaymentStatus   string    `json:"payment_status" gorm:"type:varchar(50);"` // Payment status (e.g., Pending, Completed, Failed)
	CallbackURL     string    `json:"callback_url" gorm:"type:varchar(255);"`  // Callback URL for payment notifications
	CustomerPhone   string    `json:"customer_phone" gorm:"type:varchar(20);"` // Customer's phone number
	CustomerName    string    `json:"customer_name" gorm:"type:varchar(100);"` // Customer's name
	AccountReference string   `json:"account_reference" gorm:"type:varchar(100);"` // Account reference (e.g., order ID)
	TransactionDesc string    `json:"transaction_desc" gorm:"type:varchar(255);"` // Transaction description	
	TransactionDate string	  `json:"transaction_date" gorm:"type:varchar(255);"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;"` // Timestamp when the payment was created
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;"` // Timestamp when the payment was last updated
	// Relationships
}

//ratings
type Rating struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primary_key"`
	SellerID   uuid.UUID      `json:"seller_id" gorm:"type:varchar(36);"` // Reference to seller
    UserID    uuid.UUID      `json:"user_id" gorm:"type:varchar(36);"`   // Who left the rating
    Stars      int            `json:"stars" gorm:"not null;check:stars>=1 AND stars<=5"`
    Comment    string         `json:"comment" gorm:"type:text"`                  // Optional feedback
	User       User           `json:"user" gorm:"foreignKey:UserID;references:ID"`
    
    // If you also want to include seller details:
    Seller     User           `json:"seller" gorm:"foreignKey:SellerID;references:ID"`
    CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

