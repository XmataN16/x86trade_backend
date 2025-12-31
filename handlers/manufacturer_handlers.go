package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"x86trade_backend/models"
	"x86trade_backend/repository"
)

// AdminGetManufacturers — получить список производителей с пагинацией
func AdminGetManufacturers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры пагинации
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

	// Вычисляем offset
	offset := (page - 1) * limit

	// Получаем производителей
	ms, err := repository.GetManufacturersWithPagination(r.Context(), limit, offset)
	if err != nil {
		log.Printf("AdminGetManufacturers: repository error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Получаем общее количество производителей
	total, err := repository.CountManufacturers(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с пагинацией
	response := map[string]interface{}{
		"data":  ms,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminCreateManufacturer — создать производителя.
// JSON: { "name":"...","country":"...","website":"..." }
func AdminCreateManufacturer(w http.ResponseWriter, r *http.Request) {
	var payload models.Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateManufacturer(r.Context(), &payload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// AdminUpdateManufacturer — обновление производителя.
func AdminUpdateManufacturer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	var payload models.Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	payload.ID = id
	if err := repository.UpdateManufacturer(r.Context(), &payload); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminDeleteManufacturer — удалить производителя.
func AdminDeleteManufacturer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteManufacturer(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
