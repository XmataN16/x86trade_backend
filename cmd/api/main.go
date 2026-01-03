package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"x86trade_backend/internal/db"
	"x86trade_backend/internal/middleware"
	"x86trade_backend/internal/routes"
	"x86trade_backend/internal/server"

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

	// Конфиг CORS из env
	originsStr := os.Getenv("FRONTEND_ORIGIN")
	origins := []string{}
	if originsStr != "" {
		for _, o := range strings.Split(originsStr, ",") {
			if s := strings.TrimSpace(o); s != "" {
				origins = append(origins, s)
			}
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:5500", "http://127.0.0.1:5500"}
	}

	allowCred := false
	if os.Getenv("CORS_ALLOW_CREDENTIALS") == "1" || strings.ToLower(os.Getenv("CORS_ALLOW_CREDENTIALS")) == "true" {
		allowCred = true
	}

	// DEBUG
	debug := middleware.GetDebugFromEnv()
	if debug {
		log.Println("DEBUG mode enabled — each endpoint result will be logged")
	}

	// DB connect
	db.Connect()
	defer db.DB.Close()

	// Создаем роутер
	router := chi.NewRouter()

	// Применяем CORS (вынесено в internal/server)
	server.ApplyCORS(router, origins, allowCred)

	// Подключаем logging middleware (включается когда DEBUG=true)
	router.Use(middleware.LoggingMiddleware(debug))

	// Настраиваем остальные роуты
	routes.SetupRoutes(router)

	addr := ":" + os.Getenv("APP_PORT")
	if os.Getenv("APP_PORT") == "" {
		addr = ":8080"
	}

	log.Printf("Server listening on %s (CORS origins=%v allowCred=%v)\n", addr, origins, allowCred)
	log.Fatal(http.ListenAndServe(addr, router))
}
