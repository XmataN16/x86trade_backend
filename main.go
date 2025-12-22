package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"x86trade_backend/db"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env — relying on environment variables")
	} else {
		log.Println(".env loaded")
	}

	// Конфиг CORS
	originsStr := os.Getenv("FRONTEND_ORIGIN")
	origins := strings.Split(originsStr, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	if len(origins) == 0 || origins[0] == "" {
		origins = []string{"http://localhost:5500", "http://127.0.0.1:5500"}
	}

	allowCred := false
	if os.Getenv("CORS_ALLOW_CREDENTIALS") == "1" || os.Getenv("CORS_ALLOW_CREDENTIALS") == "true" {
		allowCred = true
	}

	// DB connect
	db.Connect()
	defer db.DB.Close()

	// Создаем роутер
	router := chi.NewRouter()

	// Настраиваем маршруты с CORS
	SetupRoutes(router, origins, allowCred)

	addr := ":8080"
	if p := os.Getenv("APP_PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("Server listening on %s (CORS origins=%v allowCred=%v)\n", addr, origins, allowCred)
	log.Fatal(http.ListenAndServe(addr, router))
}
