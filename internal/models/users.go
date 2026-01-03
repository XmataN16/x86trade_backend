package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	MidName      string    `json:"mid_name,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	PasswordHash string    `json:"-"`
}
