package models

import "database/sql"

type DeliveryMethod struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Description   sql.NullString  `json:"description,omitempty"`
	BaseCost      float64         `json:"base_cost"`
	FreeThreshold sql.NullFloat64 `json:"free_threshold,omitempty"`
	EstimatedDays sql.NullInt64   `json:"estimated_days,omitempty"`
}
