package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

// ProductFilter расширён: теперь можно фильтровать по имени категории / производителя
type ProductFilter struct {
	CategoryID       *int
	CategoryName     *string
	ManufacturerID   *int
	ManufacturerName *string
	MinPrice         *float64
	MaxPrice         *float64
	Q                *string // search in name/description
	Limit            int
	Offset           int
}

func GetProducts(ctx context.Context, f *ProductFilter) ([]models.Product, error) {
	base := `SELECT p.id, p.name, p.description, p.price, p.category_id, COALESCE(c.name,''), p.manufacturer_id, COALESCE(m.name,''), p.image_path, p.stock_quantity, p.sku, p.created_at, p.updated_at
             FROM products p
             LEFT JOIN categories c ON p.category_id = c.id
             LEFT JOIN manufacturers m ON p.manufacturer_id = m.id`
	conds := []string{}
	args := []interface{}{}
	i := 1

	if f != nil {
		if f.CategoryID != nil {
			conds = append(conds, fmt.Sprintf("p.category_id=$%d", i))
			args = append(args, *f.CategoryID)
			i++
		}
		if f.CategoryName != nil && strings.TrimSpace(*f.CategoryName) != "" {
			conds = append(conds, fmt.Sprintf("c.name ILIKE $%d", i))
			args = append(args, "%"+*f.CategoryName+"%")
			i++
		}
		if f.ManufacturerID != nil {
			conds = append(conds, fmt.Sprintf("p.manufacturer_id=$%d", i))
			args = append(args, *f.ManufacturerID)
			i++
		}
		if f.ManufacturerName != nil && strings.TrimSpace(*f.ManufacturerName) != "" {
			conds = append(conds, fmt.Sprintf("m.name ILIKE $%d", i))
			args = append(args, "%"+*f.ManufacturerName+"%")
			i++
		}
		if f.MinPrice != nil {
			conds = append(conds, fmt.Sprintf("p.price >= $%d", i))
			args = append(args, *f.MinPrice)
			i++
		}
		if f.MaxPrice != nil {
			conds = append(conds, fmt.Sprintf("p.price <= $%d", i))
			args = append(args, *f.MaxPrice)
			i++
		}
		if f.Q != nil && strings.TrimSpace(*f.Q) != "" {
			conds = append(conds, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", i, i+1))
			args = append(args, "%"+*f.Q+"%", "%"+*f.Q+"%")
			i += 2
		}
	}

	if len(conds) > 0 {
		base = base + " WHERE " + strings.Join(conds, " AND ")
	}

	base = base + " ORDER BY p.id DESC"

	if f != nil && f.Limit > 0 {
		base = base + fmt.Sprintf(" LIMIT %d", f.Limit)
	}
	if f != nil && f.Offset > 0 {
		base = base + fmt.Sprintf(" OFFSET %d", f.Offset)
	}

	rows, err := db.DB.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Product
	for rows.Next() {
		var p models.Product
		var created, updated sql.NullTime
		var categoryName, manufacturerName sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &categoryName, &p.ManufacturerID, &manufacturerName, &p.ImagePath, &p.StockQuantity, &p.SKU, &created, &updated); err != nil {
			return nil, err
		}
		if categoryName.Valid {
			p.CategoryName = categoryName.String
		}
		if manufacturerName.Valid {
			p.ManufacturerName = manufacturerName.String
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
	row := db.DB.QueryRowContext(ctx, `SELECT p.id, p.name, p.description, p.price, p.category_id, COALESCE(c.name,''), p.manufacturer_id, COALESCE(m.name,''), p.image_path, p.stock_quantity, p.sku, p.created_at, p.updated_at
        FROM products p
        LEFT JOIN categories c ON p.category_id=c.id
        LEFT JOIN manufacturers m ON p.manufacturer_id=m.id
        WHERE p.id=$1`, id)
	var created, updated sql.NullTime
	var categoryName, manufacturerName sql.NullString
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &categoryName, &p.ManufacturerID, &manufacturerName, &p.ImagePath, &p.StockQuantity, &p.SKU, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if categoryName.Valid {
		p.CategoryName = categoryName.String
	}
	if manufacturerName.Valid {
		p.ManufacturerName = manufacturerName.String
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
	err := db.DB.QueryRowContext(ctx, `INSERT INTO products (name, description, price, category_id, manufacturer_id, image_path, stock_quantity, sku) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		p.Name, p.Description, p.Price, p.CategoryID, p.ManufacturerID, p.ImagePath, p.StockQuantity, p.SKU).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateProduct(ctx context.Context, p *models.Product) error {
	_, err := db.DB.ExecContext(ctx, `UPDATE products SET name=$1, description=$2, price=$3, category_id=$4, manufacturer_id=$5, image_path=$6, stock_quantity=$7, sku=$8, updated_at = CURRENT_TIMESTAMP WHERE id=$9`,
		p.Name, p.Description, p.Price, p.CategoryID, p.ManufacturerID, p.ImagePath, p.StockQuantity, p.SKU, p.ID)
	return err
}

func DeleteProduct(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}
