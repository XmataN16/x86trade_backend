package repository

import (
	"context"
	"time"
	"x86trade_backend/db"
	"x86trade_backend/models"
)

func CreateReview(ctx context.Context, review *models.Review) error {
	_, err := db.DB.ExecContext(ctx, `
		INSERT INTO reviews (product_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, review.ProductID, review.UserID, review.Rating, review.Comment, time.Now())

	return err
}
