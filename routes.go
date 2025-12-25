package main

import (
	"x86trade_backend/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// SetupRoutes теперь принимает готовый роутер и конфигурацию CORS
func SetupRoutes(r chi.Router, origins []string, allowCred bool) {
	// Сначала добавляем CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: allowCred,
		MaxAge:           300,
		// Debug: true, // раскомментировать для отладки CORS
	}))

	// Публичные роуты
	r.Post("/api/auth/register", handlers.RegisterHandler)
	r.Post("/api/auth/login", handlers.LoginHandler)
	r.Post("/api/auth/refresh", handlers.RefreshHandler)
	r.Post("/api/auth/logout", handlers.LogoutHandler)

	// Обратная связь
	r.Post("/api/contact", handlers.CreateContactMessageHandler)

	// Вакании
	r.Get("/api/vacancies", handlers.GetVacanciesHandler)

	// Продукты
	r.Get("/api/products", handlers.GetProductsHandler)
	r.Get("/api/products/{id}", handlers.GetProductHandler)

	// Категории
	r.Get("/api/categories", handlers.GetCategoriesHandler)

	// Методы доставки и оплаты
	r.Get("/api/delivery_methods", handlers.GetDeliveryMethodsHandler)
	r.Get("/api/payment_methods", handlers.GetPaymentMethodsHandler)

	// Характеристики товара с указанным id
	r.Get("/api/products/{id}/details", handlers.GetProductDetailsHandler)

	// Защищенные роуты (требуют аутентификации)
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		// Корзина
		r.Get("/api/cart", handlers.GetCartHandler)
		r.Post("/api/cart", handlers.AddToCartHandler)
		r.Put("/api/cart", handlers.UpdateCartHandler)
		r.Delete("/api/cart", handlers.ClearCartHandler)
		r.Delete("/api/cart/{productID}", handlers.RemoveFromCartHandler)

		// Заказы
		r.Get("/api/orders", handlers.GetOrdersHandler)
		r.Get("/api/orders/{id}", handlers.GetOrderDetailsHandler) // Обновленный роут для получения деталей с товарами
		r.Post("/api/orders", handlers.CreateOrderHandler)
		r.Put("/api/orders/{orderID}/cancel", handlers.CancelOrderHandler)

		// Профиль
		r.Get("/api/auth/me", handlers.MeHandler)
		r.Put("/api/auth/me", handlers.UpdateMeHandler)

		r.Post("/api/reviews", handlers.CreateReviewHandler)
	})
}
