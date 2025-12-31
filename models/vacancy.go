package models

import "time"

type Vacancy struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Requirements string    `json:"requirements"`
	Conditions   string    `json:"conditions"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type VacancyPayload struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Requirements string `json:"requirements"`
	Conditions   string `json:"conditions"`
	ContactEmail string `json:"contact_email"`
}
