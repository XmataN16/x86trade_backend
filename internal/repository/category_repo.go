package repository

import (
	"context"
	"database/sql"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

// GetCategories возвращает все категории (без вложений).
func GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, slug, image_path FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		var img sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Slug, &img); err != nil {
			return nil, err
		}
		if img.Valid {
			c.ImagePath = img.String
		}
		out = append(out, c)
	}
	return out, nil
}

// CreateCategory вставляет категорию и возвращает её id.
func CreateCategory(ctx context.Context, c *models.Category) (int, error) {
	q := `INSERT INTO categories (name, slug) VALUES ($1,$2) RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, q, c.Name, nullableString(c.Slug)).Scan(&id)
	return id, err
}

// GetCategoryByID возвращает категорию по id (nil, nil если нет).
func GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	q := `SELECT id, name, slug FROM categories WHERE id=$1`
	var c models.Category
	var slug sql.NullString
	if err := db.DB.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if slug.Valid {
		c.Slug = slug.String
	}
	return &c, nil
}

// GetAllCategories возвращает все категории (без пагинации).
func GetAllCategories(ctx context.Context) ([]models.Category, error) {
	q := `SELECT id, name, description, slug FROM categories ORDER BY id`
	rows, err := db.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		var descr sql.NullString
		var slug sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &slug); err != nil {
			return nil, err
		}
		if slug.Valid {
			c.Slug = slug.String
		}
		if descr.Valid {
			c.Description = descr.String
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpdateCategory обновляет категорию по id.
func UpdateCategory(ctx context.Context, c *models.Category) error {
	q := `UPDATE categories SET name=$1, slug=$2 WHERE id=$3`
	_, err := db.DB.ExecContext(ctx, q, c.Name, nullableString(c.Slug), c.ID)
	return err
}

// DeleteCategory удаляет категорию.
func DeleteCategory(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM categories WHERE id=$1`, id)
	return err
}
