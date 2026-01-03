package routes

import (
	"x86trade_backend/internal/handlers"
	"x86trade_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes регистрирует публичные и защищённые пользовательские роуты.
// CORS и logging теперь применяются извне (в main).
func SetupRoutes(r chi.Router) {
	// Публичные роуты
	r.Post("/api/auth/register", handlers.RegisterHandler)
	r.Post("/api/auth/login", handlers.LoginHandler)
	r.Post("/api/auth/refresh", handlers.RefreshHandler)
	r.Post("/api/auth/logout", handlers.LogoutHandler)

	r.Post("/api/contact", handlers.CreateContactMessageHandler)
	r.Get("/api/vacancies", handlers.GetVacanciesHandler)

	r.Get("/api/products", handlers.GetProductsHandler)
	r.Get("/api/products/{id}", handlers.GetProductHandler)
	r.Get("/api/products/{id}/details", handlers.GetProductDetailsHandler)

	r.Get("/api/categories", handlers.GetCategoriesHandler)
	r.Get("/api/delivery_methods", handlers.GetDeliveryMethodsHandler)
	r.Get("/api/payment_methods", handlers.GetPaymentMethodsHandler)

	// Защищенные роуты (требуют аутентификации)
	r.Group(func(r chi.Router) {
		// Auth middleware — у вас уже реализовано в handlers.AuthMiddleware
		r.Use(middleware.AuthMiddleware)

		// Cart
		r.Get("/api/cart", handlers.GetCartHandler)
		r.Post("/api/cart", handlers.AddToCartHandler)
		r.Put("/api/cart", handlers.UpdateCartHandler)
		r.Delete("/api/cart", handlers.ClearCartHandler)
		r.Delete("/api/cart/{productID}", handlers.RemoveFromCartHandler)

		// Orders (user)
		r.Get("/api/orders", handlers.GetOrdersHandler)
		r.Get("/api/orders/{id}", handlers.GetOrderDetailsHandler)
		r.Post("/api/orders", handlers.CreateOrderHandler)
		r.Put("/api/orders/{orderID}/cancel", handlers.CancelOrderHandler)

		// Profile
		r.Get("/api/auth/me", handlers.MeHandler)
		r.Put("/api/auth/me", handlers.UpdateMeHandler)

		r.Post("/api/reviews", handlers.CreateReviewHandler)

		// Админские роуты — регистрируем в отдельном модуле
		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminOnly)
			RegisterAdminRoutes(r)
		})
	})
}
