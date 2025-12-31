package repository

import (
	"context"
	"database/sql"
	"time"
	"x86trade_backend/db"
	"x86trade_backend/models"
)

func GetVacancies(ctx context.Context) ([]models.Vacancy, error) {
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, title, description, requirements, conditions, contact_email, created_at, updated_at
		FROM vacancies
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []models.Vacancy
	for rows.Next() {
		var v models.Vacancy
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&v.ID, &v.Title, &v.Description,
			&v.Requirements, &v.Conditions, &v.ContactEmail,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			v.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			v.UpdatedAt = updatedAt.Time
		}

		vacancies = append(vacancies, v)
	}

	return vacancies, nil
}

// CreateVacancy вставляет новую вакансию и возвращает её id.
func CreateVacancy(ctx context.Context, v *models.Vacancy) (int, error) {
	q := `
		INSERT INTO vacancies
			(title, description, requirements, conditions, contact_email, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`
	now := time.Now().UTC()
	var id int
	err := db.DB.QueryRowContext(ctx, q,
		v.Title, v.Description, v.Requirements, v.Conditions, v.ContactEmail, now, now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetVacancyByID возвращает вакансию по id (nil, nil если не найдено).
func GetVacancyByID(ctx context.Context, id int) (*models.Vacancy, error) {
	q := `
		SELECT id, title, description, requirements, conditions, contact_email, created_at, updated_at
		FROM vacancies WHERE id = $1
	`
	var v models.Vacancy
	var createdAt, updatedAt sql.NullTime
	err := db.DB.QueryRowContext(ctx, q, id).Scan(
		&v.ID, &v.Title, &v.Description, &v.Requirements, &v.Conditions, &v.ContactEmail,
		&createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		v.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		v.UpdatedAt = updatedAt.Time
	}
	return &v, nil
}

// UpdateVacancy обновляет вакансию (поля title/description/... ), и обновляет updated_at.
func UpdateVacancy(ctx context.Context, v *models.Vacancy) error {
	q := `
		UPDATE vacancies SET
			title = $1,
			description = $2,
			requirements = $3,
			conditions = $4,
			contact_email = $5,
			updated_at = $6
		WHERE id = $7
	`
	_, err := db.DB.ExecContext(ctx, q,
		v.Title, v.Description, v.Requirements, v.Conditions, v.ContactEmail, time.Now().UTC(), v.ID,
	)
	return err
}

// DeleteVacancy удаляет вакансию по id.
func DeleteVacancy(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM vacancies WHERE id = $1`, id)
	return err
}
