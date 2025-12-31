package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

type ProductFilter struct {
	CategoryID       *int
	CategoryName     *string
	ManufacturerID   *int
	ManufacturerName *string
	MinPrice         *float64
	MaxPrice         *float64
	Q                *string
	Limit            int
	Offset           int
}

func GetProducts(ctx context.Context, f *ProductFilter) ([]models.Product, error) {
	base := `SELECT p.id, p.name, p.description, p.price, p.category_id, 
                    COALESCE(c.name,'') AS category_name, 
                    p.manufacturer_id, 
                    COALESCE(m.name,'') AS manufacturer_name, 
                    p.image_path, p.stock_quantity, p.sku, 
                    p.created_at, p.updated_at
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
			args = append(args, "%"+strings.TrimSpace(*f.CategoryName)+"%")
			i++
		}
		if f.ManufacturerID != nil {
			conds = append(conds, fmt.Sprintf("p.manufacturer_id=$%d", i))
			args = append(args, *f.ManufacturerID)
			i++
		}
		if f.ManufacturerName != nil && strings.TrimSpace(*f.ManufacturerName) != "" {
			conds = append(conds, fmt.Sprintf("m.name ILIKE $%d", i))
			args = append(args, "%"+strings.TrimSpace(*f.ManufacturerName)+"%")
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
			search := "%" + strings.TrimSpace(*f.Q) + "%"
			conds = append(conds, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", i, i+1))
			args = append(args, search, search)
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

	// Добавляем логирование для отладки
	log.Printf("Executing query: %s", base)
	log.Printf("With args: %v", args)

	rows, err := db.DB.QueryContext(ctx, base, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var out []models.Product
	for rows.Next() {
		var p models.Product
		var created, updated sql.NullTime
		var categoryName, manufacturerName, description, imagePath, sku sql.NullString
		var stockQuantity sql.NullInt64

		err := rows.Scan(&p.ID, &p.Name, &description, &p.Price, &p.CategoryID,
			&categoryName, &p.ManufacturerID, &manufacturerName,
			&imagePath, &stockQuantity, &sku, &created, &updated)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		// Обработка NULL значений
		if description.Valid {
			p.Description = description.String
		}
		if categoryName.Valid {
			p.CategoryName = categoryName.String
		}
		if manufacturerName.Valid {
			p.ManufacturerName = manufacturerName.String
		}
		if imagePath.Valid {
			p.ImagePath = imagePath.String
		}
		if sku.Valid {
			p.SKU = sku.String
		}
		if stockQuantity.Valid {
			p.StockQuantity = int(stockQuantity.Int64)
		}
		if created.Valid {
			p.CreatedAt = created.Time
		}
		if updated.Valid {
			p.UpdatedAt = updated.Time
		}
		out = append(out, p)
	}

	// Проверяем ошибки после цикла
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return out, nil
}

// GetProductByID возвращает продукт по id (nil, nil если не найден).
func GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	q := `SELECT id, name, sku, description, price, category_id, manufacturer_id, image_path, stock_quantity, created_at, updated_at 
          FROM products WHERE id=$1`
	var p models.Product
	var description, imagePath, sku sql.NullString
	var categoryID, manufacturerID, stockQuantity sql.NullInt64
	var created, updated sql.NullTime

	err := db.DB.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.Name, &sku, &description, &p.Price,
		&categoryID, &manufacturerID, &imagePath, &stockQuantity, &created, &updated)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Обработка NULL значений
	if description.Valid {
		p.Description = description.String
	}
	if categoryID.Valid {
		p.CategoryID = int(categoryID.Int64)
	}
	if manufacturerID.Valid {
		p.ManufacturerID = int(manufacturerID.Int64)
	}
	if imagePath.Valid {
		p.ImagePath = imagePath.String
	}
	if stockQuantity.Valid {
		p.StockQuantity = int(stockQuantity.Int64)
	}
	if sku.Valid {
		p.SKU = sku.String
	}
	if created.Valid {
		p.CreatedAt = created.Time
	}
	if updated.Valid {
		p.UpdatedAt = updated.Time
	}

	return &p, nil
}

// GetAllProducts возвращает все продукты (без пагинации). Можно позже расширить limit/offset.
func GetAllProducts(ctx context.Context) ([]models.Product, error) {
	q := `SELECT id, name, sku, description, price, category_id, manufacturer_id, image_path, stock_quantity, created_at FROM products ORDER BY id`
	rows, err := db.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Product
	for rows.Next() {
		var p models.Product
		var categoryID sql.NullInt64
		var manufacturerID sql.NullInt64
		var imgPath sql.NullString
		var stock_quantity sql.NullInt64
		var created sql.NullTime
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Description, &p.Price, &categoryID, &manufacturerID, &imgPath, &stock_quantity, &created); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			p.CategoryID = int(categoryID.Int64)
		}
		if manufacturerID.Valid {
			p.ManufacturerID = int(manufacturerID.Int64)
		}
		if imgPath.Valid {
			p.ImagePath = imgPath.String
		}
		if stock_quantity.Valid {
			p.StockQuantity = int(stock_quantity.Int64)
		}
		if created.Valid {
			p.CreatedAt = created.Time
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// CreateProduct вставляет новый продукт и возвращает id.
func CreateProduct(ctx context.Context, p *models.Product) (int, error) {
	q := `INSERT INTO products (name, sku, description, price, category_id, image_path, stock_quantity, manufacturer_id, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`
	var id int
	now := time.Now().UTC()
	err := db.DB.QueryRowContext(ctx, q,
		p.Name, p.SKU, p.Description, p.Price, nullableInt(p.CategoryID), p.ImagePath, p.StockQuantity, nullableInt(p.ManufacturerID), now).Scan(&id)
	return id, err
}

// UpdateProduct обновляет поля продукта по id.
func UpdateProduct(ctx context.Context, p *models.Product) error {
	q := `UPDATE products SET name=$1, sku=$2, description=$3, price=$4, category_id=$5, manufacturer_id=$6, image_path=$7, stock_quantity=$8 WHERE id=$9`
	_, err := db.DB.ExecContext(ctx, q, p.Name, p.SKU, p.Description, p.Price, nullableInt(p.CategoryID), nullableInt(p.ManufacturerID), p.ImagePath, p.StockQuantity, p.ID)
	return err
}

// DeleteProduct удаляет продукт по id.
func DeleteProduct(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

// helper: чтобы передавать NULL для 0 (если в вашей модели 0 означает пусто)
func nullableInt(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func GetProductDetails(ctx context.Context, productID int) (*models.ProductDetail, error) {
	// Получаем основную информацию о товаре
	product, err := GetProductByID(ctx, productID)
	if err != nil {
		log.Printf("Error getting product %d: %v", productID, err)
		return nil, err
	}
	if product == nil {
		log.Printf("Product with ID %d not found", productID)
		return nil, nil
	}

	log.Printf("Found product: %s (ID: %d)", product.Name, productID)

	// Получаем характеристики товара с обработкой ошибок
	characteristics, err := GetProductCharacteristics(ctx, productID)
	if err != nil {
		log.Printf("Error getting characteristics (continuing without them): %v", err)
		characteristics = []models.ProductCharacteristic{}
	}

	log.Printf("Characteristics count: %d", len(characteristics))

	// Получаем отзывы о товаре с обработкой ошибок
	reviews, err := GetProductReviews(ctx, productID)
	if err != nil {
		log.Printf("Error getting reviews (continuing without them): %v", err)
		reviews = []models.Review{}
	}

	log.Printf("Reviews count: %d", len(reviews))

	// Вычисляем средний рейтинг
	var averageRating float64
	if len(reviews) > 0 {
		var totalRating int
		for _, r := range reviews {
			totalRating += r.Rating
		}
		averageRating = float64(totalRating) / float64(len(reviews))
		log.Printf("Average rating calculated: %.2f", averageRating)
	}

	detail := &models.ProductDetail{
		Product:         *product,
		Characteristics: characteristics,
		Reviews:         reviews,
		AverageRating:   averageRating,
	}

	return detail, nil
}

func GetProductCharacteristics(ctx context.Context, productID int) ([]models.ProductCharacteristic, error) {
	rows, err := db.DB.QueryContext(ctx, `
        SELECT pc.id, pc.product_id, ct.name as characteristic_name, ct.unit as characteristic_unit, pc.value
        FROM product_characteristics pc
        LEFT JOIN characteristic_types ct ON pc.characteristic_type_id = ct.id
        WHERE pc.product_id = $1
        ORDER BY ct.name
    `, productID)
	if err != nil {
		log.Printf("Error getting characteristics for product %d: %v", productID, err)
		return nil, err
	}
	defer rows.Close()

	var characteristics []models.ProductCharacteristic
	for rows.Next() {
		var c models.ProductCharacteristic
		err := rows.Scan(&c.ID, &c.ProductID, &c.CharacteristicName, &c.CharacteristicUnit, &c.Value)
		if err != nil {
			log.Printf("Error scanning characteristic row: %v", err)
			continue // Продолжаем даже при ошибке сканирования одной строки
		}
		characteristics = append(characteristics, c)
	}

	// Проверяем ошибки после цикла
	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning characteristics: %v", err)
		return nil, err
	}

	log.Printf("Found %d characteristics for product %d", len(characteristics), productID)
	return characteristics, nil
}

func GetProductReviews(ctx context.Context, productID int) ([]models.Review, error) {
	rows, err := db.DB.QueryContext(ctx, `
        SELECT r.id, r.product_id, r.user_id, u.first_name || ' ' || u.last_name as user_name, 
               r.rating, r.comment, r.created_at
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        WHERE r.product_id = $1
        ORDER BY r.created_at DESC
    `, productID)
	if err != nil {
		log.Printf("Error getting reviews for product %d: %v", productID, err)
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		var userName sql.NullString
		var createdAt sql.NullTime

		if err := rows.Scan(&r.ID, &r.ProductID, &r.UserID, &userName, &r.Rating, &r.Comment, &createdAt); err != nil {
			log.Printf("Error scanning review row: %v", err)
			return nil, err
		}

		if userName.Valid {
			r.UserName = userName.String
		} else {
			r.UserName = "Аноним"
		}

		if createdAt.Valid {
			r.CreatedAt = createdAt.Time
		}

		reviews = append(reviews, r)
	}

	// Проверка ошибки после цикла
	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning reviews: %v", err)
		return nil, err
	}

	log.Printf("Found %d reviews for product %d", len(reviews), productID)
	return reviews, nil
}
