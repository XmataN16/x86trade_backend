package admin_handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"
)

// AdminGetPaymentMethods — возвращаем список (переиспользуем GetPaymentMethods).
func AdminGetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	methods, err := repository.GetPaymentMethods(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(methods)
}

// AdminCreatePaymentMethod — принимает простой JSON, конвертирует в repository.PaymentMethod.
func AdminCreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	var p models.PaymentPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	var pm models.PaymentMethod
	pm.Name = p.Name
	if p.Description != nil {
		pm.Description = sql.NullString{String: *p.Description, Valid: true}
	} else {
		pm.Description = sql.NullString{}
	}
	// По умолчанию считаем активным, если не указан
	if p.IsActive != nil {
		pm.IsActive = *p.IsActive
	} else {
		pm.IsActive = true
	}

	id, err := repository.CreatePaymentMethod(r.Context(), &pm)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// AdminUpdatePaymentMethod — ожидаем полный payload (name + is_active желательно).
func AdminUpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}

	var p models.PaymentPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	// Требуем явный is_active в update, чтобы не нечаянно сбросить флаг.
	if p.IsActive == nil {
		http.Error(w, "is_active required", http.StatusBadRequest)
		return
	}

	var pm models.PaymentMethod
	pm.ID = id
	pm.Name = p.Name
	if p.Description != nil {
		pm.Description = sql.NullString{String: *p.Description, Valid: true}
	} else {
		pm.Description = sql.NullString{}
	}
	pm.IsActive = *p.IsActive

	if err := repository.UpdatePaymentMethod(r.Context(), &pm); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminDeletePaymentMethod — удаление по id.
func AdminDeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeletePaymentMethod(r.Context(), id); err != nil {
		// если FK в orders -> вернётся ошибка; даём понятный код при конфликте
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
