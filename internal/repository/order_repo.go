package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

// Helper: get cart items for user (reuse existing GetCartByUserID)...
// Создание заказа в транзакции.
func CreateOrderFromCart(ctx context.Context, userID int, deliveryMethodID *int, address, recipientName, recipientPhone, comment string) (int, error) {
	// 1) достать корзину
	cartItems, err := GetCartByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if len(cartItems) == 0 {
		return 0, errors.New("cart empty")
	}

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// посчитать total, получить текущую цену товара (защита от изменения)
	total := 0.0
	// получим цену для каждого товара
	for i := range cartItems {
		var price sql.NullFloat64
		err = tx.QueryRowContext(ctx, `SELECT price FROM products WHERE id=$1`, cartItems[i].ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		p := 0.0
		if price.Valid {
			p = price.Float64
		}
		total += p * float64(cartItems[i].Quantity)
	}

	// вставляем заказ
	var orderID int
	now := time.Now()
	// status 'created'
	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (user_id, status, total_amount, created_at, updated_at, comment) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		userID, "created", total, now, now, comment).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// вставляем order_items
	for _, ci := range cartItems {
		var price sql.NullFloat64
		if err = tx.QueryRowContext(ctx, `SELECT price FROM products WHERE id=$1`, ci.ProductID).Scan(&price); err != nil {
			tx.Rollback()
			return 0, err
		}
		p := 0.0
		if price.Valid {
			p = price.Float64
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO order_items (order_id, product_id, quantity, price_per_unit) VALUES ($1,$2,$3,$4)`,
			orderID, ci.ProductID, ci.Quantity, p)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// вставка данных доставки (если переданы)
	if deliveryMethodID != nil && address != "" {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_deliveries (order_id, delivery_method_id, address, recipient_name, recipient_phone, status) VALUES ($1,$2,$3,$4,$5,$6)`,
			orderID, *deliveryMethodID, address, recipientName, recipientPhone, "pending")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Очистка корзины пользователя
	if _, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id=$1`, userID); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return orderID, nil
}

// Получение заказов пользователя (простой вариант)
func GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, user_id, status, total_amount, created_at, updated_at, comment FROM orders WHERE user_id=$1 ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Order
	for rows.Next() {
		var o models.Order
		var created, updated sql.NullTime
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &created, &updated, &o.Comment); err != nil {
			return nil, err
		}
		if created.Valid {
			o.CreatedAt = created.Time
		}
		if updated.Valid {
			o.UpdatedAt = updated.Time
		}
		out = append(out, o)
	}
	return out, nil
}

// Получение деталей заказа (включая позиции)
func GetOrderWithItems(ctx context.Context, orderID int) (*models.Order, []models.OrderItem, error) {
	var ord models.Order
	var created, updated sql.NullTime
	row := db.DB.QueryRowContext(ctx, `SELECT id, user_id, status, total_amount, created_at, updated_at, comment FROM orders WHERE id=$1`, orderID)
	if err := row.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.TotalAmount, &created, &updated, &ord.Comment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if created.Valid {
		ord.CreatedAt = created.Time
	}
	if updated.Valid {
		ord.UpdatedAt = updated.Time
	}

	rows, err := db.DB.QueryContext(ctx, `SELECT id, order_id, product_id, quantity, price_per_unit, total_price FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return &ord, nil, err
	}
	defer rows.Close()
	var items []models.OrderItem
	for rows.Next() {
		var it models.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.PricePerUnit, &it.TotalPrice); err != nil {
			return &ord, nil, err
		}
		items = append(items, it)
	}
	return &ord, items, nil
}

func GetOrderDelivery(ctx context.Context, orderID int) (*models.OrderDelivery, error) {
	var delivery models.OrderDelivery
	row := db.DB.QueryRowContext(ctx, `
        SELECT od.address, od.recipient_name, od.recipient_phone, od.status,
               dm.name as method_name, dm.base_cost
        FROM order_deliveries od
        JOIN delivery_methods dm ON od.delivery_method_id = dm.id
        WHERE od.order_id = $1
    `, orderID)

	err := row.Scan(&delivery.Address, &delivery.RecipientName,
		&delivery.RecipientPhone, &delivery.Status,
		&delivery.MethodName, &delivery.BaseCost)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &delivery, nil
}

// GetOrdersWithPagination возвращает список заказов с пагинацией
func GetOrdersWithPagination(ctx context.Context, limit, offset int) ([]models.Order, error) {
	q := `SELECT id, user_id, status, total_amount, created_at, updated_at, comment 
		  FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := db.DB.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var created, updated sql.NullTime

		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &created, &updated, &o.Comment); err != nil {
			return nil, err
		}

		if created.Valid {
			o.CreatedAt = created.Time
		}
		if updated.Valid {
			o.UpdatedAt = updated.Time
		}

		orders = append(orders, o)
	}

	return orders, nil
}

// CountOrders возвращает общее количество заказов
func CountOrders(ctx context.Context) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&count)
	return count, err
}

// UpdateOrderStatus обновляет статус заказа
func UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	_, err := db.DB.ExecContext(ctx, `
		UPDATE orders 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2
	`, status, orderID)

	return err
}
