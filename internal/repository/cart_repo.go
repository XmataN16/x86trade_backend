package repository

import (
	"context"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

// GetCartByUserID returns cart items for given user id
func GetCartByUserID(ctx context.Context, userID int) ([]models.CartItem, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, user_id, product_id, quantity FROM cart_items WHERE user_id=$1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CartItem
	for rows.Next() {
		var it models.CartItem
		if err := rows.Scan(&it.ID, &it.UserID, &it.ProductID, &it.Quantity); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, nil
}

func AddOrUpdateCartItem(ctx context.Context, userID int, productID int, quantity int) error {
	// Используем PostgreSQL upsert: при конфликте по (user_id, product_id) увеличим quantity.
	// Предполагается, что в таблице есть UNIQUE (user_id, product_id).
	_, err := db.DB.ExecContext(ctx, `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, userID, productID, quantity)
	return err
}

func RemoveCartItem(ctx context.Context, userID int, productID int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id=$1 AND product_id=$2`, userID, productID)
	return err
}

func ClearCart(ctx context.Context, userID int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id=$1`, userID)
	return err
}
