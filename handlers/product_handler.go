package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"x86trade_backend/models"
	"x86trade_backend/repository"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := repository.GetProducts(ctx)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	p, err := repository.GetProductByID(context.Background(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if p == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateProduct(context.Background(), &p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}
