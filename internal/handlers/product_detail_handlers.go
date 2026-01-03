package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"x86trade_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func GetProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		http.Error(w, "bad request: invalid product id", http.StatusBadRequest)
		return
	}

	log.Printf("Getting details for product ID: %d", productID)

	productDetail, err := repository.GetProductDetails(r.Context(), productID)
	if err != nil {
		log.Printf("Error getting product details: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if productDetail == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productDetail)
}
