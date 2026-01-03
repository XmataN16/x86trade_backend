package admin_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func AdminGetOrders(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	orders, err := repository.GetOrdersWithPagination(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ordersWithDetails := make([]map[string]interface{}, 0, len(orders))
	for _, order := range orders {
		_, items, err := repository.GetOrderWithItems(r.Context(), order.ID)
		if err != nil {
			continue
		}

		var deliveryInfo *models.OrderDelivery
		deliveryInfo, _ = repository.GetOrderDelivery(r.Context(), order.ID)

		user, _ := repository.GetUserByID(r.Context(), order.UserID)

		orderMap := map[string]interface{}{
			"id":           order.ID,
			"user_id":      order.UserID,
			"user_name":    "",
			"status":       order.Status,
			"total_amount": order.TotalAmount,
			"created_at":   order.CreatedAt.Format("2006-01-02 15:04"),
			"updated_at":   order.UpdatedAt.Format("2006-01-02 15:04"),
			"comment":      order.Comment,
			"items":        items,
		}

		if user != nil {
			orderMap["user_name"] = user.FirstName + " " + user.LastName
		}

		if deliveryInfo != nil {
			orderMap["delivery"] = map[string]interface{}{
				"method_name":     deliveryInfo.MethodName,
				"address":         deliveryInfo.Address,
				"recipient_name":  deliveryInfo.RecipientName,
				"recipient_phone": deliveryInfo.RecipientPhone,
				"status":          deliveryInfo.Status,
			}
		}

		ordersWithDetails = append(ordersWithDetails, orderMap)
	}

	total, err := repository.CountOrders(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  ordersWithDetails,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AdminGetOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		http.Error(w, "bad request: invalid order id", http.StatusBadRequest)
		return
	}

	ord, items, err := repository.GetOrderWithItems(r.Context(), orderID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	deliveryInfo, _ := repository.GetOrderDelivery(r.Context(), orderID)

	user, _ := repository.GetUserByID(r.Context(), ord.UserID)

	response := map[string]interface{}{
		"order":    ord,
		"items":    items,
		"delivery": deliveryInfo,
		"user":     user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AdminUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		http.Error(w, "bad request: invalid order id", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	if payload.Status == "" {
		http.Error(w, "bad request: status is required", http.StatusBadRequest)
		return
	}

	ord, _, err := repository.GetOrderWithItems(r.Context(), orderID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := repository.UpdateOrderStatus(r.Context(), orderID, payload.Status); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order status updated successfully",
	})
}

func AdminUpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	if payload.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	if err := repository.UpdateOrderStatus(r.Context(), orderID, payload.Status); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
