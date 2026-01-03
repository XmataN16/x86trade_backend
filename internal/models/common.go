package models

type DeliveryPayload struct {
	Name          string   `json:"name"`
	Description   *string  `json:"description,omitempty"`
	BaseCost      float64  `json:"base_cost"`
	FreeThreshold *float64 `json:"free_threshold,omitempty"`
	EstimatedDays *int64   `json:"estimated_days,omitempty"`
}

type CreateOrderPayload struct {
	DeliveryMethodID *int   `json:"delivery_method_id,omitempty"`
	Address          string `json:"address,omitempty"`
	RecipientName    string `json:"recipient_name,omitempty"`
	RecipientPhone   string `json:"recipient_phone,omitempty"`
	Comment          string `json:"comment,omitempty"`
}
