package repository

import (
	"context"
	"database/sql"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

// GetCategories возвращает все категории (без вложений).
func GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, parent_id, slug, image_path FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		var parent sql.NullInt64
		var img sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &parent, &c.Slug, &img); err != nil {
			return nil, err
		}
		if parent.Valid {
			v := int(parent.Int64)
			c.ParentID = &v
		} else {
			c.ParentID = nil
		}
		if img.Valid {
			c.ImagePath = img.String
		}
		out = append(out, c)
	}
	return out, nil
}

func GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	var c models.Category
	var parent sql.NullInt64
	var img sql.NullString
	row := db.DB.QueryRowContext(ctx, `SELECT id, name, description, parent_id, slug, image_path FROM categories WHERE id=$1`, id)
	if err := row.Scan(&c.ID, &c.Name, &c.Description, &parent, &c.Slug, &img); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if parent.Valid {
		v := int(parent.Int64)
		c.ParentID = &v
	}
	if img.Valid {
		c.ImagePath = img.String
	}
	return &c, nil
}

func CreateCategory(ctx context.Context, c *models.Category) (int, error) {
	var id int
	// parent_id может быть nil; image_path тоже может быть nil/пустой
	if c.ParentID != nil {
		err := db.DB.QueryRowContext(ctx,
			`INSERT INTO categories (name, description, parent_id, slug, image_path) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
			c.Name, c.Description, c.ParentID, c.Slug, c.ImagePath).Scan(&id)
		if err != nil {
			return 0, err
		}
	} else {
		err := db.DB.QueryRowContext(ctx,
			`INSERT INTO categories (name, description, slug, image_path) VALUES ($1,$2,$3,$4) RETURNING id`,
			c.Name, c.Description, c.Slug, c.ImagePath).Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func UpdateCategory(ctx context.Context, c *models.Category) error {
	_, err := db.DB.ExecContext(ctx, `UPDATE categories SET name=$1, description=$2, parent_id=$3, slug=$4, image_path=$5 WHERE id=$6`,
		c.Name, c.Description, c.ParentID, c.Slug, c.ImagePath, c.ID)
	return err
}

func DeleteCategory(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM categories WHERE id=$1`, id)
	return err
}
