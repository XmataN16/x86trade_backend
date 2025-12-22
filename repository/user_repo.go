package repository

import (
	"context"
	"database/sql"
	"time"

	"x86trade_backend/db"
	"x86trade_backend/models"
)

func CreateUser(ctx context.Context, u *models.User, passwordHash string) (int, error) {
	var id int
	err := db.DB.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, first_name, last_name) VALUES ($1,$2,$3,$4) RETURNING id`,
		u.Email, passwordHash, u.FirstName, u.LastName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	row := db.DB.QueryRowContext(ctx, `SELECT id, email, first_name, last_name, phone, is_admin, created_at, password_hash FROM users WHERE email=$1`, email)
	var created sql.NullTime
	var phone sql.NullString
	var pass sql.NullString
	if err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &phone, &u.IsAdmin, &created, &pass); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if phone.Valid {
		u.Phone = phone.String
	}
	if created.Valid {
		u.CreatedAt = created.Time
	}
	if pass.Valid {
		u.PasswordHash = pass.String
	}
	return &u, nil
}

func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	row := db.DB.QueryRowContext(ctx, `SELECT id, email, first_name, last_name, phone, is_admin, created_at, password_hash FROM users WHERE id=$1`, id)
	var created sql.NullTime
	var phone sql.NullString
	var pass sql.NullString
	if err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &phone, &u.IsAdmin, &created, &pass); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if phone.Valid {
		u.Phone = phone.String
	}
	if created.Valid {
		u.CreatedAt = created.Time
	}
	if pass.Valid {
		u.PasswordHash = pass.String
	}
	return &u, nil
}

// Refresh tokens ops
func SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	_, err := db.DB.ExecContext(ctx, `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1,$2,$3)`, userID, token, expiresAt)
	return err
}

func DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token=$1`, token)
	return err
}

func GetRefreshToken(ctx context.Context, token string) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time
	err := db.DB.QueryRowContext(ctx, `SELECT user_id, expires_at FROM refresh_tokens WHERE token=$1`, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, time.Time{}, err
	}
	return userID, expiresAt, nil
}

// UpdateUserProfile обновляет поля профиля (first_name, last_name, phone) для пользователя userID.
func UpdateUserProfile(ctx context.Context, userID int, firstName, lastName, phone string) error {
	_, err := db.DB.ExecContext(ctx,
		`UPDATE users SET first_name=$1, last_name=$2, phone=$3 WHERE id=$4`,
		firstName, lastName, phone, userID,
	)
	return err
}
