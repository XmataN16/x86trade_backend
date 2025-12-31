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

		// Админская зона — требует аутентификацию и админские права
		r.Group(func(r chi.Router) {
			r.Use(handlers.AdminOnly)

			// users CRUD (admin)
			r.Get("/api/admin/users", handlers.AdminGetUsers)
			r.Get("/api/admin/users/{id}", handlers.AdminGetUserByID)
			r.Post("/api/admin/users", handlers.AdminCreateUser)
			r.Put("/api/admin/users/{id}", handlers.AdminUpdateUser)
			r.Delete("/api/admin/users/{id}", handlers.AdminDeleteUser)

			// products CRUD (admin)
			r.Get("/api/admin/products", handlers.AdminGetProducts)
			r.Get("/api/admin/products/{id}", handlers.AdminGetProductByID)
			r.Post("/api/admin/products", handlers.AdminCreateProduct)
			r.Put("/api/admin/products/{id}", handlers.AdminUpdateProduct)
			r.Delete("/api/admin/products/{id}", handlers.AdminDeleteProduct)

			// categories CRUD (admin)
			r.Get("/api/admin/categories", handlers.AdminGetCategories)
			r.Post("/api/admin/categories", handlers.AdminCreateCategory)
			r.Put("/api/admin/categories/{id}", handlers.AdminUpdateCategory)
			r.Delete("/api/admin/categories/{id}", handlers.AdminDeleteCategory)

			// manufacturers CRUD (admin)
			r.Get("/api/admin/manufacturers", handlers.AdminGetManufacturers)
			r.Post("/api/admin/manufacturers", handlers.AdminCreateManufacturer)
			r.Put("/api/admin/manufacturers/{id}", handlers.AdminUpdateManufacturer)
			r.Delete("/api/admin/manufacturers/{id}", handlers.AdminDeleteManufacturer)

			// delivery_methods CRUD (admin)
			r.Get("/api/admin/delivery_methods", handlers.AdminGetDeliveryMethods)
			r.Post("/api/admin/delivery_methods", handlers.AdminCreateDeliveryMethod)
			r.Put("/api/admin/delivery_methods/{id}", handlers.AdminUpdateDeliveryMethod)
			r.Delete("/api/admin/delivery_methods/{id}", handlers.AdminDeleteDeliveryMethod)

			// payment methods CRUD (admin)
			r.Get("/api/admin/payment_methods", handlers.AdminGetPaymentMethods)
			r.Post("/api/admin/payment_methods", handlers.AdminCreatePaymentMethod)
			r.Put("/api/admin/payment_methods/{id}", handlers.AdminUpdatePaymentMethod)
			r.Delete("/api/admin/payment_methods/{id}", handlers.AdminDeletePaymentMethod)

			// vacancies CRUD (admin)
			r.Get("/api/admin/vacancies", handlers.AdminGetVacancies)
			r.Get("/api/admin/vacancies/{id}", handlers.AdminGetVacancyByID)
			r.Post("/api/admin/vacancies", handlers.AdminCreateVacancy)
			r.Put("/api/admin/vacancies/{id}", handlers.AdminUpdateVacancy)
			r.Delete("/api/admin/vacancies/{id}", handlers.AdminDeleteVacancy)

			// characteristic types
			r.Get("/api/admin/characteristic_types", handlers.AdminGetCharacteristicTypes)
			r.Post("/api/admin/characteristic_types", handlers.AdminCreateCharacteristicType)
			r.Put("/api/admin/characteristic_types/{id}", handlers.AdminUpdateCharacteristicType)
			r.Delete("/api/admin/characteristic_types/{id}", handlers.AdminDeleteCharacteristicType)

			// product characteristics
			r.Get("/api/admin/product_characteristics", handlers.AdminListProductCharacteristics)
			r.Get("/api/admin/products/{product_id}/characteristics", handlers.AdminGetProductCharacteristics)
			r.Post("/api/admin/product_characteristics", handlers.AdminCreateProductCharacteristic)
			r.Put("/api/admin/product_characteristics/{id}", handlers.AdminUpdateProductCharacteristic)
			r.Delete("/api/admin/product_characteristics/{id}", handlers.AdminDeleteProductCharacteristic)

			// orders CRUD (admin)
			r.Get("/api/admin/orders", handlers.AdminGetOrders)
			r.Put("/api/admin/orders/{id}/status", handlers.AdminUpdateOrderStatus)
			r.Put("/api/admin/orders/{id}", handlers.AdminUpdateOrder)

			// replace-all (bulk) for product
			r.Put("/api/admin/products/{product_id}/characteristics", handlers.AdminReplaceProductCharacteristics)

		})
	})
}
