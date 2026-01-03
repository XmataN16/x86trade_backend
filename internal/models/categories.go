package models

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
	ImagePath   string `json:"image_path,omitempty"`
}
