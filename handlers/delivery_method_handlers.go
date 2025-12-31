package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"x86trade_backend/models"
	"x86trade_backend/repository"
)

// AdminGetDeliveryMethods — возвращаем через существующий репозиторий.
// Используем repository.GetDeliveryMethods (у вас уже есть) чтобы поведение оставалось единым.
func AdminGetDeliveryMethods(w http.ResponseWriter, r *http.Request) {
	methods, err := repository.GetDeliveryMethods(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(methods)
}

// Входной payload: удобные простые поля, nullable поля как опциональные (указатели).
type deliveryPayload struct {
	Name          string   `json:"name"`
	Description   *string  `json:"description,omitempty"`
	BaseCost      float64  `json:"base_cost"`
	FreeThreshold *float64 `json:"free_threshold,omitempty"`
	EstimatedDays *int64   `json:"estimated_days,omitempty"`
}

// AdminCreateDeliveryMethod — принимает простой JSON, конвертирует в models.DeliveryMethod и вызывает Create.
func AdminCreateDeliveryMethod(w http.ResponseWriter, r *http.Request) {
	var p deliveryPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	var d models.DeliveryMethod
	d.Name = p.Name
	if p.Description != nil {
		d.Description = sql.NullString{String: *p.Description, Valid: true}
	} else {
		d.Description = sql.NullString{}
	}
	d.BaseCost = p.BaseCost
	if p.FreeThreshold != nil {
		d.FreeThreshold = sql.NullFloat64{Float64: *p.FreeThreshold, Valid: true}
	} else {
		d.FreeThreshold = sql.NullFloat64{}
	}
	if p.EstimatedDays != nil {
		d.EstimatedDays = sql.NullInt64{Int64: *p.EstimatedDays, Valid: true}
	} else {
		d.EstimatedDays = sql.NullInt64{}
	}

	id, err := repository.CreateDeliveryMethod(r.Context(), &d)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// AdminUpdateDeliveryMethod
func AdminUpdateDeliveryMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	var p deliveryPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	var d models.DeliveryMethod
	d.ID = id
	d.Name = p.Name
	if p.Description != nil {
		d.Description = sql.NullString{String: *p.Description, Valid: true}
	} else {
		d.Description = sql.NullString{}
	}
	d.BaseCost = p.BaseCost
	if p.FreeThreshold != nil {
		d.FreeThreshold = sql.NullFloat64{Float64: *p.FreeThreshold, Valid: true}
	} else {
		d.FreeThreshold = sql.NullFloat64{}
	}
	if p.EstimatedDays != nil {
		d.EstimatedDays = sql.NullInt64{Int64: *p.EstimatedDays, Valid: true}
	} else {
		d.EstimatedDays = sql.NullInt64{}
	}

	if err := repository.UpdateDeliveryMethod(r.Context(), &d); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminDeleteDeliveryMethod
func AdminDeleteDeliveryMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteDeliveryMethod(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
