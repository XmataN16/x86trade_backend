package repository

import (
	"context"
	"database/sql"
	"time"
	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

// CreateManufacturer вставляет запись и возвращает id.
func CreateManufacturer(ctx context.Context, m *models.Manufacturer) (int, error) {
	q := `INSERT INTO manufacturers (name, country, website, created_at) VALUES ($1,$2,$3,$4) RETURNING id`
	var id int
	now := time.Now().UTC()
	err := db.DB.QueryRowContext(ctx, q,
		m.Name, nullableString(m.Country), nullableString(m.Website), now).Scan(&id)
	return id, err
}

// GetManufacturerByID возвращает производителя по id.
func GetManufacturerByID(ctx context.Context, id int) (*models.Manufacturer, error) {
	q := `SELECT id, name, country, website, created_at FROM manufacturers WHERE id=$1`
	var m models.Manufacturer
	var country sql.NullString
	var website sql.NullString
	var created sql.NullTime
	err := db.DB.QueryRowContext(ctx, q, id).Scan(&m.ID, &m.Name, &country, &website, &created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if country.Valid {
		m.Country = country.String
	}
	if website.Valid {
		m.Website = website.String
	}
	if created.Valid {
		m.CreatedAt = created.Time
	}
	return &m, nil
}

// GetAllManufacturers возвращает всех производителей.
func GetAllManufacturers(ctx context.Context) ([]models.Manufacturer, error) {
	q := `SELECT id, name, country, website, created_at FROM manufacturers ORDER BY id`
	rows, err := db.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Manufacturer
	for rows.Next() {
		var m models.Manufacturer
		var country sql.NullString
		var website sql.NullString
		var created sql.NullTime
		if err := rows.Scan(&m.ID, &m.Name, &country, &website, &created); err != nil {
			return nil, err
		}
		if country.Valid {
			m.Country = country.String
		}
		if website.Valid {
			m.Website = website.String
		}
		if created.Valid {
			m.CreatedAt = created.Time
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// GetManufacturersWithPagination возвращает список производителей с пагинацией
func GetManufacturersWithPagination(ctx context.Context, limit, offset int) ([]models.Manufacturer, error) {
	q := `SELECT id, name, country, website, created_at FROM manufacturers ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := db.DB.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Manufacturer
	for rows.Next() {
		var m models.Manufacturer
		var country sql.NullString
		var website sql.NullString
		var created sql.NullTime
		if err := rows.Scan(&m.ID, &m.Name, &country, &website, &created); err != nil {
			return nil, err
		}
		if country.Valid {
			m.Country = country.String
		}
		if website.Valid {
			m.Website = website.String
		}
		if created.Valid {
			m.CreatedAt = created.Time
		}
		out = append(out, m)
	}

	return out, rows.Err()
}

// CountManufacturers возвращает общее количество производителей
func CountManufacturers(ctx context.Context) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM manufacturers`).Scan(&count)
	return count, err
}

// UpdateManufacturer обновляет производителя по id.
func UpdateManufacturer(ctx context.Context, m *models.Manufacturer) error {
	q := `UPDATE manufacturers SET name=$1, country=$2, website=$3 WHERE id=$4`
	_, err := db.DB.ExecContext(ctx, q, m.Name, nullableString(m.Country), nullableString(m.Website), m.ID)
	return err
}

// DeleteManufacturer удаляет производителя по id.
func DeleteManufacturer(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM manufacturers WHERE id=$1`, id)
	return err
}

// вспомогательная функция: пустая строка -> NULL
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
