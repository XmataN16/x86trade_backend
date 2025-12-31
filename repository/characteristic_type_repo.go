package repository

import (
	"context"
	"database/sql"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

func CreateCharacteristicType(ctx context.Context, t *models.CharacteristicType) (int, error) {
	q := `INSERT INTO characteristic_types (name, unit, category_id) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, q, t.Name, nullableString(t.Unit), t.CategoryID).Scan(&id)
	return id, err
}

func GetCharacteristicTypeByID(ctx context.Context, id int) (*models.CharacteristicType, error) {
	q := `SELECT id, name, unit, category_id FROM characteristic_types WHERE id=$1`
	var t models.CharacteristicType
	var unit sql.NullString
	var categoryID sql.NullInt64
	err := db.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Name, &unit, &categoryID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if unit.Valid {
		t.Unit = unit.String
	}
	if categoryID.Valid {
		t.CategoryID = int(categoryID.Int64)
	}
	return &t, nil
}

func GetAllCharacteristicTypes(ctx context.Context) ([]models.CharacteristicType, error) {
	q := `SELECT id, name, unit, category_id FROM characteristic_types ORDER BY name`
	rows, err := db.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CharacteristicType
	for rows.Next() {
		var t models.CharacteristicType
		var unit sql.NullString
		var categoryID sql.NullInt64
		if err := rows.Scan(&t.ID, &t.Name, &unit, &categoryID); err != nil {
			return nil, err
		}
		if unit.Valid {
			t.Unit = unit.String
		}
		if categoryID.Valid {
			t.CategoryID = int(categoryID.Int64)
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func UpdateCharacteristicType(ctx context.Context, t *models.CharacteristicType) error {
	q := `UPDATE characteristic_types SET name=$1, unit=$2, category_id=$3 WHERE id=$4`
	_, err := db.DB.ExecContext(ctx, q, t.Name, nullableString(t.Unit), t.CategoryID, t.ID)
	return err
}

func DeleteCharacteristicType(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM characteristic_types WHERE id=$1`, id)
	return err
}
