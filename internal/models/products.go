package models

import "time"

type Product struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	Price            float64   `json:"price"`
	CategoryID       int       `json:"category_id,omitempty"`
	CategoryName     string    `json:"category_name,omitempty"`
	ManufacturerID   int       `json:"manufacturer_id,omitempty"`
	ManufacturerName string    `json:"manufacturer_name,omitempty"`
	ImagePath        string    `json:"image_path"`
	StockQuantity    int       `json:"stock_quantity,omitempty"`
	SKU              string    `json:"sku,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

type ProductDetail struct {
	Product         Product                 `json:"product"`
	Characteristics []ProductCharacteristic `json:"characteristics"`
	Reviews         []Review                `json:"reviews"`
	AverageRating   float64                 `json:"average_rating"`
}
