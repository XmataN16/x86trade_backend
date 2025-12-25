package models

import "time"

type Product struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	Price            float64   `json:"price"`
	CategoryID       int       `json:"category_id,omitempty"`   // можно оставить для создания/обновления
	CategoryName     string    `json:"category_name,omitempty"` // NEW: имя категории для фронта
	ManufacturerID   int       `json:"manufacturer_id,omitempty"`
	ManufacturerName string    `json:"manufacturer_name,omitempty"` // NEW: имя производителя для фронта
	ImagePath        string    `json:"image_path,omitempty"`
	StockQuantity    int       `json:"stock_quantity"`
	SKU              string    `json:"sku,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"`
	Slug        string `json:"slug,omitempty"`
	ImagePath   string `json:"image_path,omitempty"`
}

// Для примера: User с минимальными полями
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	// PasswordHash is intentionally omitted from JSON
	PasswordHash string `json:"-"`
}

type CartItem struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price_per_unit,omitempty"` // optional
}

type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	Comment     string    `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type OrderItem struct {
	ID           int     `json:"id"`
	OrderID      int     `json:"order_id"`
	ProductID    int     `json:"product_id"`
	Quantity     int     `json:"quantity"`
	PricePerUnit float64 `json:"price_per_unit"`
	TotalPrice   float64 `json:"total_price"`
	// Optional: include product name/image for convenience when returning
	ProductName string `json:"product_name,omitempty"`
	ImagePath   string `json:"image_path,omitempty"`
}

type OrderDelivery struct {
	Address        string  `json:"address"`
	RecipientName  string  `json:"recipient_name"`
	RecipientPhone string  `json:"recipient_phone"`
	Status         string  `json:"status"`
	MethodName     string  `json:"method_name"`
	BaseCost       float64 `json:"base_cost"`
}

type ContactMessage struct {
	ID              int        `json:"id"`
	FullName        string     `json:"full_name"`
	ContactInfo     string     `json:"contact_info"`
	Message         string     `json:"message"`
	CreatedAt       time.Time  `json:"created_at"`
	IsProcessed     bool       `json:"is_processed"`
	ResponseMessage string     `json:"response_message,omitempty"`
	ResponseAt      *time.Time `json:"response_at,omitempty"`
}

type Vacancy struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Requirements string    `json:"requirements"`
	Conditions   string    `json:"conditions"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductDetail struct {
	Product         Product                 `json:"product"`
	Characteristics []ProductCharacteristic `json:"characteristics"`
	Reviews         []Review                `json:"reviews"`
	AverageRating   float64                 `json:"average_rating"`
}
