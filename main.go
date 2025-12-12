package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"x86trade_backend/db"
	"x86trade_backend/handlers"
)

// withCORS оборачивает handler и добавляет заголовки CORS.
// origin может быть "*" или конкретный origin "http://localhost:5500".
func withCORS(h http.Handler, origin string, allowCredentials bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем Access-Control-Allow-Origin.
		// Если origin == "*", то используем "*". Иначе указываем конкретный origin.
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if allowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// Примечание: при allowCredentials=true нельзя использовать "*" в Access-Control-Allow-Origin.
		}

		// Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	// загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env — relying on environment variables")
	} else {
		log.Println(".env loaded")
	}

	// Конфиг CORS
	// Поставь FRONTEND_ORIGIN в .env, например: FRONTEND_ORIGIN=http://localhost:5500
	origin := os.Getenv("FRONTEND_ORIGIN")
	// Для разработки удобно разрешать все, но если планируешь использовать cookie/auth, укажи точный origin и set allowCredentials=true
	allowCred := false
	if os.Getenv("CORS_ALLOW_CREDENTIALS") == "1" || os.Getenv("CORS_ALLOW_CREDENTIALS") == "true" {
		allowCred = true
	}

	// DB connect
	db.Connect()
	defer db.DB.Close()

	// Используем ServeMux, затем оборачиваем entire mux в CORS
	mux := http.NewServeMux()

	// auth routes
	mux.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/auth/login", handlers.LoginHandler)
	mux.HandleFunc("/api/auth/refresh", handlers.RefreshHandler)
	mux.HandleFunc("/api/auth/logout", handlers.LogoutHandler)

	// cart routes (protected) — оборачиваем AuthMiddleware
	mux.Handle("/api/cart", handlers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetCartHandler(w, r)
		case http.MethodPost:
			handlers.AddToCartHandler(w, r)
		case http.MethodPut:
			handlers.UpdateCartHandler(w, r)
		case http.MethodDelete:
			// if query param product_id present -> remove single, else clear
			if r.URL.Query().Get("product_id") != "" {
				handlers.RemoveFromCartHandler(w, r)
			} else {
				handlers.ClearCartHandler(w, r)
			}
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Products: GET, POST, PUT, DELETE
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetProductsHandler(w, r)
		case http.MethodPost:
			handlers.CreateProductHandler(w, r)
		case http.MethodPut:
			handlers.UpdateProductHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteProductHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Categories
	mux.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				handlers.GetCategoryHandler(w, r)
			} else {
				handlers.GetCategoriesHandler(w, r)
			}
		case http.MethodPost:
			handlers.CreateCategoryHandler(w, r)
		case http.MethodPut:
			handlers.UpdateCategoryHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteCategoryHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":8080"
	if p := os.Getenv("APP_PORT"); p != "" {
		addr = ":" + p
	}

	handlerWithCORS := withCORS(mux, origin, allowCred)

	log.Printf("Server listening on %s (CORS origin=%s allowCred=%v)\n", addr, origin, allowCred)
	log.Fatal(http.ListenAndServe(addr, handlerWithCORS))
}
