package repository

import (
	"context"
	"database/sql"

	"x86trade_backend/db"
)

type PaymentMethod struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description,omitempty"`
	IsActive    bool           `json:"is_active"`
}

func GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, is_active FROM payment_methods ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PaymentMethod
	for rows.Next() {
		var p PaymentMethod
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.IsActive); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}
