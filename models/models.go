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
}
