package repository

import (
	"context"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

func GetPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, is_active FROM payment_methods ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.PaymentMethod
	for rows.Next() {
		var p models.PaymentMethod
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.IsActive); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// CreatePaymentMethod вставляет метод оплаты и возвращает id.
func CreatePaymentMethod(ctx context.Context, p *models.PaymentMethod) (int, error) {
	q := `INSERT INTO payment_methods (name, description, is_active) VALUES ($1,$2,$3) RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, q, p.Name, p.Description, p.IsActive).Scan(&id)
	return id, err
}

// UpdatePaymentMethod обновляет метод оплаты по id.
func UpdatePaymentMethod(ctx context.Context, p *models.PaymentMethod) error {
	q := `UPDATE payment_methods SET name=$1, description=$2, is_active=$3 WHERE id=$4`
	_, err := db.DB.ExecContext(ctx, q, p.Name, p.Description, p.IsActive, p.ID)
	return err
}

// DeletePaymentMethod удаляет метод оплаты по id.
func DeletePaymentMethod(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM payment_methods WHERE id=$1`, id)
	return err
}
