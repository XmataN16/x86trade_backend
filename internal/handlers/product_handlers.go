package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"
)

func parseIntPtr(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseFloatPtr(s string) (*float64, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// если указан id — вернуть единичный ресурс
	if idStr := r.URL.Query().Get("id"); idStr != "" {
		GetProductHandler(w, r)
		return
	}

	q := r.URL.Query()
	filter := &repository.ProductFilter{}

	if cid := q.Get("category_id"); cid != "" {
		if v, err := strconv.Atoi(cid); err == nil {
			filter.CategoryID = &v
		}
	} else if cname := q.Get("category"); cname != "" {
		cname = strings.TrimSpace(cname)
		if cname != "" {
			filter.CategoryName = &cname
		}
	}

	if brand := q.Get("brand"); brand != "" {
		b := strings.TrimSpace(brand)
		if b != "" {
			filter.ManufacturerName = &b
		}
	}

	if min := q.Get("min_price"); min != "" {
		if f, err := strconv.ParseFloat(min, 64); err == nil {
			filter.MinPrice = &f
		}
	}
	if max := q.Get("max_price"); max != "" {
		if f, err := strconv.ParseFloat(max, 64); err == nil {
			filter.MaxPrice = &f
		}
	}
	if search := q.Get("q"); search != "" {
		s := strings.TrimSpace(search)
		if s != "" {
			filter.Q = &s
		}
	}
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			filter.Limit = v
		}
	}
	if o := q.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			filter.Offset = v
		}
	}

	products, err := repository.GetProducts(r.Context(), filter)
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

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}
	if p.ID == 0 {
		http.Error(w, "bad request: id required", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateProduct(context.Background(), &p); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteProduct(context.Background(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
