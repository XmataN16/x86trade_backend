package routes

import (
	"x86trade_backend/internal/handlers/admin_handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterAdminRoutes регистрирует все админские endpoint'ы.
// Вызывается внутри защищённой группы, где уже подключён middleware AdminOnly.
func RegisterAdminRoutes(r chi.Router) {
	// users CRUD (admin)
	r.Get("/api/admin/users", admin_handlers.AdminGetUsers)
	r.Get("/api/admin/users/{id}", admin_handlers.AdminGetUserByID)
	r.Post("/api/admin/users", admin_handlers.AdminCreateUser)
	r.Put("/api/admin/users/{id}", admin_handlers.AdminUpdateUser)
	r.Delete("/api/admin/users/{id}", admin_handlers.AdminDeleteUser)

	// products CRUD (admin)
	r.Get("/api/admin/products", admin_handlers.AdminGetProducts)
	r.Get("/api/admin/products/{id}", admin_handlers.AdminGetProductByID)
	r.Post("/api/admin/products", admin_handlers.AdminCreateProduct)
	r.Put("/api/admin/products/{id}", admin_handlers.AdminUpdateProduct)
	r.Delete("/api/admin/products/{id}", admin_handlers.AdminDeleteProduct)

	// categories CRUD (admin)
	r.Get("/api/admin/categories", admin_handlers.AdminGetCategories)
	r.Post("/api/admin/categories", admin_handlers.AdminCreateCategory)
	r.Put("/api/admin/categories/{id}", admin_handlers.AdminUpdateCategory)
	r.Delete("/api/admin/categories/{id}", admin_handlers.AdminDeleteCategory)

	// manufacturers CRUD (admin)
	r.Get("/api/admin/manufacturers", admin_handlers.AdminGetManufacturers)
	r.Post("/api/admin/manufacturers", admin_handlers.AdminCreateManufacturer)
	r.Put("/api/admin/manufacturers/{id}", admin_handlers.AdminUpdateManufacturer)
	r.Delete("/api/admin/manufacturers/{id}", admin_handlers.AdminDeleteManufacturer)

	// delivery_methods CRUD (admin)
	r.Get("/api/admin/delivery_methods", admin_handlers.AdminGetDeliveryMethods)
	r.Post("/api/admin/delivery_methods", admin_handlers.AdminCreateDeliveryMethod)
	r.Put("/api/admin/delivery_methods/{id}", admin_handlers.AdminUpdateDeliveryMethod)
	r.Delete("/api/admin/delivery_methods/{id}", admin_handlers.AdminDeleteDeliveryMethod)

	// payment methods CRUD (admin)
	r.Get("/api/admin/payment_methods", admin_handlers.AdminGetPaymentMethods)
	r.Post("/api/admin/payment_methods", admin_handlers.AdminCreatePaymentMethod)
	r.Put("/api/admin/payment_methods/{id}", admin_handlers.AdminUpdatePaymentMethod)
	r.Delete("/api/admin/payment_methods/{id}", admin_handlers.AdminDeletePaymentMethod)

	// vacancies CRUD (admin)
	r.Get("/api/admin/vacancies", admin_handlers.AdminGetVacancies)
	r.Get("/api/admin/vacancies/{id}", admin_handlers.AdminGetVacancyByID)
	r.Post("/api/admin/vacancies", admin_handlers.AdminCreateVacancy)
	r.Put("/api/admin/vacancies/{id}", admin_handlers.AdminUpdateVacancy)
	r.Delete("/api/admin/vacancies/{id}", admin_handlers.AdminDeleteVacancy)

	// characteristic types
	r.Get("/api/admin/characteristic_types", admin_handlers.AdminGetCharacteristicTypes)
	r.Post("/api/admin/characteristic_types", admin_handlers.AdminCreateCharacteristicType)
	r.Put("/api/admin/characteristic_types/{id}", admin_handlers.AdminUpdateCharacteristicType)
	r.Delete("/api/admin/characteristic_types/{id}", admin_handlers.AdminDeleteCharacteristicType)

	// product characteristics
	r.Get("/api/admin/product_characteristics", admin_handlers.AdminListProductCharacteristics)
	r.Get("/api/admin/products/{product_id}/characteristics", admin_handlers.AdminGetProductCharacteristics)
	r.Post("/api/admin/product_characteristics", admin_handlers.AdminCreateProductCharacteristic)
	r.Put("/api/admin/product_characteristics/{id}", admin_handlers.AdminUpdateProductCharacteristic)
	r.Delete("/api/admin/product_characteristics/{id}", admin_handlers.AdminDeleteProductCharacteristic)

	// orders CRUD (admin)
	r.Get("/api/admin/orders", admin_handlers.AdminGetOrders)
	r.Put("/api/admin/orders/{id}/status", admin_handlers.AdminUpdateOrderStatus)
	r.Put("/api/admin/orders/{id}", admin_handlers.AdminUpdateOrder)

	// replace-all (bulk) for product
	r.Put("/api/admin/products/{product_id}/characteristics", admin_handlers.AdminReplaceProductCharacteristics)
}
