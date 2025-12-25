// handlers/product_detail_handler.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"x86trade_backend/models"
	"x86trade_backend/repository"

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

func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID int    `json:"product_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	if req.ProductID <= 0 {
		http.Error(w, "bad request: invalid product id", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "bad request: rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	if len(req.Comment) < 10 {
		http.Error(w, "bad request: comment must be at least 10 characters", http.StatusBadRequest)
		return
	}

	review := &models.Review{
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := repository.CreateReview(r.Context(), review); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Review created successfully",
	})
}
