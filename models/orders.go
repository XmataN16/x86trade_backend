package models

import "time"

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
	ProductName  string  `json:"product_name,omitempty"`
	ImagePath    string  `json:"image_path,omitempty"`
}

type OrderDelivery struct {
	Address        string  `json:"address"`
	RecipientName  string  `json:"recipient_name"`
	RecipientPhone string  `json:"recipient_phone"`
	Status         string  `json:"status"`
	MethodName     string  `json:"method_name"`
	BaseCost       float64 `json:"base_cost"`
}
