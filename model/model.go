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
    Order       []Order `gorm:"foreignKey:UserID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
}


type Product struct {
    BaseModel
    Name        string  `json:"name" gorm:"size:255"`
    Description string  `json:"description" gorm:"type:text"`
    Price       float64 `json:"price" gorm:"type:decimal(10,2)"`
   	Category    string  `json:"category" gorm:"size:100"` // Example: Shirts, Trousers, etc.
    Stock       int     `json:"stock"`                  // Available stock count
    ImageURL    string  `json:"image_url" gorm:"size:255"` // URL for the clothing item image
    IsActive    bool    `json:"is_active" gorm:"default:true"` // Active status for display
    OrderItems  []OrderItem `gorm:"foreignKey:ProductID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
    CartItems   []CartItem  `gorm:"foreignKey:ProductID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
}

type Order struct {
	BaseModel
	OrderNumber   string          `json:"order_number" gorm:"size:100;unique;not null"`
	UserID        uuid.UUID       `json:"user_id" gorm:"type:varchar(36);not null"`
	User          User            `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
	TotalAmount   float64         `json:"total_amount" gorm:"type:decimal(10,2)"`
	PaymentStatus string          `json:"payment_status" gorm:"size:50"`
	PaymentMethod string          `json:"payment_method" gorm:"size:50"`
	ShippingAddress string        `json:"shipping_address" gorm:"type:text"`
	OrderStatus   string          `json:"order_status" gorm:"size:50"`
	DeliveredAt   *time.Time      `json:"delivered_at"`
	Items         []OrderItem     `json:"items" gorm:"foreignKey:OrderID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
}

type OrderItem struct {
	BaseModel
	OrderID     uuid.UUID  `json:"order_id" gorm:"type:varchar(36);"`
	Order       Order      `json:"order" gorm:"foreignKey:OrderID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	ProductID   uuid.UUID  `json:"product_id" gorm:"type:varchar(36);"`
	Product     Product    `json:"product" gorm:"foreignKey:ProductID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
	Quantity    int        `json:"quantity" gorm:"int"`
	Price       float64    `json:"price" gorm:"type:decimal(10,2)"`
	TotalPrice  float64    `json:"total_price" gorm:"type:decimal(10,2)"`
}

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
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PaymentMethod   string     `json:"payment_method" gorm:"type:varchar(100)"`  // e.g., M-Pesa, PayPal, Credit Card
	Amount          float64    `json:"amount" gorm:"type:decimal(10,2)"`         // The payment amount
	Status          string     `json:"status" gorm:"type:varchar(100)"`          // Status of the payment (e.g., Pending, Completed, Failed)
	TransactionID   string     `json:"transaction_id" gorm:"type:varchar(100)"`  // Unique transaction ID from the payment gateway
	PhoneNumber     string     `json:"phone_number" gorm:"type:varchar(15)"`     // Payer's phone number (especially for mobile payments like M-Pesa)
	PaymentDate     time.Time  `json:"payment_date" gorm:"type:timestamp"`       // Date and time of the payment
	OrderID         uuid.UUID  `json:"order_id" gorm:"type:uuid"`                // Associated order ID
	MerchantRequestID string   `json:"merchant_request_id" gorm:"type:varchar(100)"` // Request ID from M-Pesa (optional)
	ResultCode      int        `json:"result_code" gorm:"type:int"`              // Result code from payment provider (e.g., 0 for success)
	ResultMessage   string     `json:"result_message" gorm:"type:varchar(255)"`  // Message or description of the transaction result
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
