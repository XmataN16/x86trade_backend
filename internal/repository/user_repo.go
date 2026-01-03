package repository

import (
	"context"
	"database/sql"
	"time"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/models"
)

func CreateUser(ctx context.Context, u *models.User, passwordHash string) (int, error) {
	var id int
	err := db.DB.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, first_name, last_name, midname, phone) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		u.Email, passwordHash, u.FirstName, u.LastName, nullableString(u.MidName), nullableString(u.Phone)).Scan(&id)
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

// GetAllUsers возвращает список пользователей (без password_hash).
func GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := db.DB.QueryContext(ctx, `SELECT id, email, first_name, last_name, midname, phone, is_admin, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.User
	for rows.Next() {
		var u models.User
		var phone, midname sql.NullString
		var created sql.NullTime
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &midname, &phone, &u.IsAdmin, &created); err != nil {
			return nil, err
		}
		if midname.Valid {
			u.MidName = midname.String
		}
		if phone.Valid {
			u.Phone = phone.String
		}
		if created.Valid {
			u.CreatedAt = created.Time
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func UpdateUser(ctx context.Context, u *models.User) error {
	_, err := db.DB.ExecContext(ctx,
		`UPDATE users SET email=$1, first_name=$2, last_name=$3, midname=$4, phone=$5, is_admin=$6 WHERE id=$7`,
		u.Email, u.FirstName, u.LastName, nullableString(u.MidName), nullableString(u.Phone), u.IsAdmin, u.ID)
	return err
}

func UpdateUserPassword(ctx context.Context, userID int, passwordHash string) error {
	_, err := db.DB.ExecContext(ctx,
		`UPDATE users SET password_hash=$1 WHERE id=$2`,
		passwordHash, userID)
	return err
}

func DeleteUser(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

// CountUsers возвращает общее количество пользователей
func CountUsers(ctx context.Context) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}
