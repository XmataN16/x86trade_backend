package repository

import (
	"context"
	"database/sql"
	"errors"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

func GetProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, price, category_id, manufacturer_id, image_path, stock_quantity, sku, created_at, updated_at FROM products ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Product
	for rows.Next() {
		var p models.Product
		var created, updated sql.NullTime
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.ManufacturerID, &p.ImagePath, &p.StockQuantity, &p.SKU, &created, &updated); err != nil {
			return nil, err
		}
		if created.Valid {
			p.CreatedAt = created.Time
		}
		if updated.Valid {
			p.UpdatedAt = updated.Time
		}
		out = append(out, p)
	}
	return out, nil
}

func GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	var p models.Product
	row := db.DB.QueryRowContext(ctx, `SELECT id, name, description, price, category_id, manufacturer_id, image_path, stock_quantity, sku, created_at, updated_at FROM products WHERE id=$1`, id)
	var created, updated sql.NullTime
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.ManufacturerID, &p.ImagePath, &p.StockQuantity, &p.SKU, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if created.Valid {
		p.CreatedAt = created.Time
	}
	if updated.Valid {
		p.UpdatedAt = updated.Time
	}
	return &p, nil
}

func CreateProduct(ctx context.Context, p *models.Product) (int, error) {
	var id int
	err := db.DB.QueryRowContext(ctx, `INSERT INTO products (name, description, price, category_id, manufacturer_id, image_path, stock_quantity, sku) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`, p.Name, p.Description, p.Price, p.CategoryID, p.ManufacturerID, p.ImagePath, p.StockQuantity, p.SKU).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
