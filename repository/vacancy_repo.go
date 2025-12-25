package repository

import (
	"context"
	"database/sql"
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
