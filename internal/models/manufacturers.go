package models

import "time"

type Manufacturer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country,omitempty"`
	Website   string    `json:"website,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
