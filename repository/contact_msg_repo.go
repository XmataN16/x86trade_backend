package repository

import (
	"context"
	"time"
	"x86trade_backend/db"
	"x86trade_backend/models"
)

func CreateContactMessage(ctx context.Context, msg *models.ContactMessage) (int, error) {
	var id int
	err := db.DB.QueryRowContext(ctx, `
		INSERT INTO contact_messages (full_name, contact_info, message, created_at, is_processed)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, msg.FullName, msg.ContactInfo, msg.Message, time.Now(), false).Scan(&id)

	return id, err
}
