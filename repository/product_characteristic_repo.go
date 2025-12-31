package repository

import (
	"context"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

// CreateProductCharacteristic вставляет одну запись и возвращает id.
func CreateProductCharacteristic(ctx context.Context, in *models.ProductCharacteristicInput) (int, error) {
	q := `INSERT INTO product_characteristics (product_id, characteristic_type_id, value) VALUES ($1,$2,$3) RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, q, in.ProductID, in.CharacteristicTypeID, in.Value).Scan(&id)
	return id, err
}

func UpdateProductCharacteristic(ctx context.Context, in *models.ProductCharacteristicInput) error {
	q := `UPDATE product_characteristics SET characteristic_type_id=$1, value=$2 WHERE id=$3`
	_, err := db.DB.ExecContext(ctx, q, in.CharacteristicTypeID, in.Value, in.ID)
	return err
}

func DeleteProductCharacteristic(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM product_characteristics WHERE id=$1`, id)
	return err
}

// ReplaceProductCharacteristics — транзакционно заменяет все характеристики товара:
// удаляет старые и вставляет новые (useful when editing product details form).
func ReplaceProductCharacteristics(ctx context.Context, productID int, inputs []models.ProductCharacteristicInput) error {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// удаляем старые
	if _, err = tx.ExecContext(ctx, `DELETE FROM product_characteristics WHERE product_id = $1`, productID); err != nil {
		return err
	}

	// вставляем новые
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO product_characteristics (product_id, characteristic_type_id, value) VALUES ($1,$2,$3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, in := range inputs {
		if _, err := stmt.ExecContext(ctx, productID, in.CharacteristicTypeID, in.Value); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func GetAllProductCharacteristicsWithPagination(ctx context.Context, limit, offset int) ([]models.ProductCharacteristic, error) {
	rows, err := db.DB.QueryContext(ctx, `
        SELECT pc.id, pc.product_id,
               COALESCE(p.name, '') as product_name,
               COALESCE(ct.name, '') as characteristic_name,
               COALESCE(ct.unit, '') as characteristic_unit,
               pc.value
        FROM product_characteristics pc
        LEFT JOIN products p ON pc.product_id = p.id
        LEFT JOIN characteristic_types ct ON pc.characteristic_type_id = ct.id
        ORDER BY pc.id
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.ProductCharacteristic
	for rows.Next() {
		var c models.ProductCharacteristic
		if err := rows.Scan(&c.ID, &c.ProductID, &c.ProductName, &c.CharacteristicName, &c.CharacteristicUnit, &c.Value); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func CountProductCharacteristics(ctx context.Context) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM product_characteristics`).Scan(&count)
	return count, err
}
