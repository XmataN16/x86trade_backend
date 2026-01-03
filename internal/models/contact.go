package models

import "time"

type ContactMessage struct {
	ID              int        `json:"id"`
	FullName        string     `json:"full_name"`
	ContactInfo     string     `json:"contact_info"`
	Message         string     `json:"message"`
	CreatedAt       time.Time  `json:"created_at"`
	IsProcessed     bool       `json:"is_processed"`
	ResponseMessage string     `json:"response_message,omitempty"`
	ResponseAt      *time.Time `json:"response_at,omitempty"`
}
