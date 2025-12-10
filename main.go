package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"x86trade_backend/db"
	"x86trade_backend/handlers"
)

func main() {
	// Попытка загрузить .env (если файла нет — продолжим и будем полагаться на системные env)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env — relying on environment variables")
	} else {
		log.Println(".env loaded")
	}

	// (не фатальное предупреждение) — теперь переменные должны быть доступны
	if os.Getenv("DB_HOST") == "" {
		log.Println("Warning: DB_HOST not set")
	}

	// Подключаем БД (в db.Connect используются дефолты, если vars пусты)
	db.Connect()
	defer db.DB.Close()

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				handlers.GetProductHandler(w, r)
			} else {
				handlers.GetProductsHandler(w, r)
			}
		case http.MethodPost:
			handlers.CreateProductHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":8080"
	if p := os.Getenv("APP_PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
