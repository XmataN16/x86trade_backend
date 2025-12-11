package repository

import (
	"context"
	"database/sql"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

// GetCategories возвращает все категории (без вложений).
func GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, name, description, parent_id, slug FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		var parent sql.NullInt64
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &parent, &c.Slug); err != nil {
			return nil, err
		}
		if parent.Valid {
			v := int(parent.Int64)
			c.ParentID = &v
		} else {
			c.ParentID = nil
		}
		out = append(out, c)
	}
	return out, nil
}

func GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	var c models.Category
	var parent sql.NullInt64
	row := db.DB.QueryRowContext(ctx, `SELECT id, name, description, parent_id, slug FROM categories WHERE id=$1`, id)
	if err := row.Scan(&c.ID, &c.Name, &c.Description, &parent, &c.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if parent.Valid {
		v := int(parent.Int64)
		c.ParentID = &v
	}
	return &c, nil
}

func CreateCategory(ctx context.Context, c *models.Category) (int, error) {
	var id int
	// parent_id может быть nil
	if c.ParentID != nil {
		err := db.DB.QueryRowContext(ctx, `INSERT INTO categories (name, description, parent_id, slug) VALUES ($1,$2,$3,$4) RETURNING id`,
			c.Name, c.Description, c.ParentID, c.Slug).Scan(&id)
		if err != nil {
			return 0, err
		}
	} else {
		err := db.DB.QueryRowContext(ctx, `INSERT INTO categories (name, description, slug) VALUES ($1,$2,$3) RETURNING id`,
			c.Name, c.Description, c.Slug).Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func UpdateCategory(ctx context.Context, c *models.Category) error {
	// Обновляем все поля (ожидается полная модель). Можно сделать patch позже.
	_, err := db.DB.ExecContext(ctx, `UPDATE categories SET name=$1, description=$2, parent_id=$3, slug=$4 WHERE id=$5`,
		c.Name, c.Description, c.ParentID, c.Slug, c.ID)
	return err
}

func DeleteCategory(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM categories WHERE id=$1`, id)
	return err
}
