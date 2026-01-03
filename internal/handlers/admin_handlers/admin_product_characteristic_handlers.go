package admin_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func AdminListProductCharacteristics(w http.ResponseWriter, r *http.Request) {
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

	list, err := repository.GetAllProductCharacteristicsWithPagination(r.Context(), limit, offset)
	if err != nil {
		// логируем в stdout/stderr (поможет диагностировать 500)
		// log.Printf("AdminListProductCharacteristics error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	total, err := repository.CountProductCharacteristics(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"data":  list,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func AdminGetProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, _ := strconv.Atoi(productIDStr)
	if productID <= 0 {
		http.Error(w, "bad request: product id", http.StatusBadRequest)
		return
	}
	chars, err := repository.GetProductCharacteristics(r.Context(), productID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(chars)
}

func AdminCreateProductCharacteristic(w http.ResponseWriter, r *http.Request) {
	var p models.ProductCharacteristicInput
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.ProductID <= 0 || p.CharacteristicTypeID <= 0 || p.Value == "" {
		http.Error(w, "product_id, characteristic_type_id and value required", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateProductCharacteristic(r.Context(), &p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func AdminUpdateProductCharacteristic(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	var p models.ProductCharacteristicInput
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p.ID = id
	if p.ID <= 0 || p.CharacteristicTypeID <= 0 || p.Value == "" {
		http.Error(w, "id, characteristic_type_id and value required", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateProductCharacteristic(r.Context(), &p); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AdminDeleteProductCharacteristic(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteProductCharacteristic(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AdminReplaceProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, _ := strconv.Atoi(productIDStr)
	if productID <= 0 {
		http.Error(w, "bad request: product id", http.StatusBadRequest)
		return
	}
	var inputs []models.ProductCharacteristicInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	for _, it := range inputs {
		if it.CharacteristicTypeID <= 0 || it.Value == "" {
			http.Error(w, "each characteristic must include characteristic_type_id and non-empty value", http.StatusBadRequest)
			return
		}
		it.ProductID = productID
	}
	if err := repository.ReplaceProductCharacteristics(r.Context(), productID, inputs); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
