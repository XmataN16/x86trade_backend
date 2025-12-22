package repository

import (
	"context"
	"database/sql"

	"x86trade_backend/db"
)

type DeliveryMethod struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Description   sql.NullString  `json:"description,omitempty"`
	BaseCost      float64         `json:"base_cost"`
	FreeThreshold sql.NullFloat64 `json:"free_threshold,omitempty"`
	EstimatedDays sql.NullInt64   `json:"estimated_days,omitempty"`
}

func GetDeliveryMethods(ctx context.Context) ([]DeliveryMethod, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, base_cost, free_threshold, estimated_days FROM delivery_methods ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DeliveryMethod
	for rows.Next() {
		var d DeliveryMethod
		if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.BaseCost, &d.FreeThreshold, &d.EstimatedDays); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}
