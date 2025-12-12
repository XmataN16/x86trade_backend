package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"x86trade_backend/repository"
)

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	items, err := repository.GetCartByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.ProductID <= 0 || payload.Quantity <= 0 {
		http.Error(w, "product_id and positive quantity required", http.StatusBadRequest)
		return
	}
	if err := repository.AddOrUpdateCartItem(r.Context(), userID, payload.ProductID, payload.Quantity); err != nil {
		// логируем ошибку в stdout/stderr для диагностики
		log.Printf("AddOrUpdateCartItem error user=%d product=%d qty=%d: %v\n", userID, payload.ProductID, payload.Quantity, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateCartHandler(w http.ResponseWriter, r *http.Request) {
	// same as AddToCart but requires id present
	AddToCartHandler(w, r)
}

func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	pidStr := r.URL.Query().Get("product_id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "product_id required", http.StatusBadRequest)
		return
	}
	if err := repository.RemoveCartItem(r.Context(), userID, pid); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := repository.ClearCart(r.Context(), userID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
