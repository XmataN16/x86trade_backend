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

type ProductCharacteristic struct {
	ID                 int    `json:"id"`
	ProductID          int    `json:"product_id"`
	ProductName        string `json:"product_name,omitempty"`
	CharacteristicName string `json:"characteristic_name"`
	CharacteristicUnit string `json:"characteristic_unit"`
	Value              string `json:"value"`
}

type CharacteristicType struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Unit       string `json:"unit,omitempty"`
	CategoryID int    `json:"category_id,omitempty"`
}

// DTO для создания/обновления значения характеристики
type ProductCharacteristicInput struct {
	ID                   int    `json:"id,omitempty"`
	ProductID            int    `json:"product_id"`
	CharacteristicTypeID int    `json:"characteristic_type_id"`
	Value                string `json:"value"`
}
