package repository

import (
	"context"
	"time"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

func GetDeliveryMethods(ctx context.Context) ([]models.DeliveryMethod, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, base_cost, free_threshold, estimated_days FROM delivery_methods ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.DeliveryMethod
	for rows.Next() {
		var d models.DeliveryMethod
		if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.BaseCost, &d.FreeThreshold, &d.EstimatedDays); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// CreateDeliveryMethod вставляет метод доставки и возвращает id.
func CreateDeliveryMethod(ctx context.Context, d *models.DeliveryMethod) (int, error) {
	q := `INSERT INTO delivery_methods
	       (name, description, base_cost, free_threshold, estimated_days, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	var id int
	now := time.Now().UTC()
	err := db.DB.QueryRowContext(ctx, q,
		d.Name, d.Description, d.BaseCost, d.FreeThreshold, d.EstimatedDays, now).Scan(&id)
	return id, err
}

// UpdateDeliveryMethod обновляет существующий метод доставки.
func UpdateDeliveryMethod(ctx context.Context, d *models.DeliveryMethod) error {
	q := `UPDATE delivery_methods SET
	        name = $1,
	        description = $2,
	        base_cost = $3,
	        free_threshold = $4,
	        estimated_days = $5
	      WHERE id = $6`
	_, err := db.DB.ExecContext(ctx, q,
		d.Name, d.Description, d.BaseCost, d.FreeThreshold, d.EstimatedDays, d.ID)
	return err
}

// DeleteDeliveryMethod удаляет метод доставки по id.
func DeleteDeliveryMethod(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM delivery_methods WHERE id=$1`, id)
	return err
}
