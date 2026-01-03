package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"x86trade_backend/internal/db"
	"x86trade_backend/internal/middleware"
	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload models.CreateOrderPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	orderID, err := repository.CreateOrderFromCart(r.Context(), userID, payload.DeliveryMethodID, payload.Address, payload.RecipientName, payload.RecipientPhone, payload.Comment)
	if err != nil {
		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"order_id": orderID})
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orders, err := repository.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	ord, items, err := repository.GetOrderWithItems(r.Context(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"order": ord, "items": items})
}

func GetOrdersWithItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем все заказы пользователя
	orders, err := repository.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Для каждого заказа получаем детали
	var ordersWithItems []map[string]interface{}
	for _, order := range orders {
		_, items, err := repository.GetOrderWithItems(r.Context(), order.ID)
		if err != nil {
			log.Printf("Error getting items for order %d: %v", order.ID, err)
			continue
		}

		// Расширяем информацию о заказе
		orderMap := map[string]interface{}{
			"id":           order.ID,
			"status":       order.Status,
			"total_amount": order.TotalAmount,
			"created_at":   order.CreatedAt.Format("2006-01-02 15:04"),
			"updated_at":   order.UpdatedAt.Format("2006-01-02 15:04"),
			"comment":      order.Comment,
			"items":        items,
		}

		// Добавляем информацию о доставке, если есть
		deliveryInfo, err := repository.GetOrderDelivery(r.Context(), order.ID)
		if err == nil && deliveryInfo != nil {
			orderMap["delivery"] = deliveryInfo
		}

		ordersWithItems = append(ordersWithItems, orderMap)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordersWithItems)
}

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderID") // Используем chi.URLParam
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		http.Error(w, "bad request: invalid order id", http.StatusBadRequest)
		return
	}

	// Получаем текущий заказ
	var ord models.Order
	row := db.DB.QueryRowContext(r.Context(), `SELECT id, user_id, status FROM orders WHERE id=$1`, orderID)
	err = row.Scan(&ord.ID, &ord.UserID, &ord.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching order: %v", err) // Логируем ошибку
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Проверяем, что пользователь может отменить этот заказ
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if ord.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Проверяем, можно ли отменить заказ (только если статус 'created' или 'processing')
	if ord.Status != "created" && ord.Status != "processing" {
		http.Error(w, "bad request: cannot cancel order in current status", http.StatusBadRequest)
		return
	}

	// Обновляем статус заказа на 'cancelled'
	_, err = db.DB.ExecContext(r.Context(), `UPDATE orders SET status='cancelled', updated_at=NOW() WHERE id=$1`, orderID)
	if err != nil {
		log.Printf("Error updating order status: %v", err) // Логируем ошибку
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Order cancelled successfully",
		"order_id":   orderID,
		"new_status": "cancelled",
	})
}

func GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		http.Error(w, "bad request: invalid order id", http.StatusBadRequest)
		return
	}

	// Получаем заказ, товары и информацию о доставке
	ord, items, err := repository.GetOrderWithItems(r.Context(), orderID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Проверяем, что пользователь имеет доступ к этому заказу
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ord.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Получаем информацию о доставке
	deliveryInfo, _ := repository.GetOrderDelivery(r.Context(), orderID)

	// Загружаем информацию о продуктах для каждого товара в заказе
	for i, item := range items {
		product, err := repository.GetProductByID(r.Context(), item.ProductID)
		if err == nil && product != nil {
			items[i].ProductName = product.Name
			items[i].ImagePath = product.ImagePath
		}
	}

	// Формируем ответ
	response := map[string]interface{}{
		"order": ord,
		"items": items,
	}

	if deliveryInfo != nil {
		response["delivery"] = deliveryInfo
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
