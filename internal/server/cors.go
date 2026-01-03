package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// ApplyCORS применяет CORS middleware к роутеру с переданной конфигурацией.
// origins — список origins, allowCred — true/false.
func ApplyCORS(r chi.Router, origins []string, allowCred bool) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: allowCred,
		MaxAge:           300,
	}))
}
