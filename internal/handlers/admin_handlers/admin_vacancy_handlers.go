package admin_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

// AdminGetVacancies — переиспользует существующую функцию GetVacancies.
func AdminGetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, err := repository.GetVacancies(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(vacancies)
}

// AdminGetVacancyByID — возвращает одну вакансию по id.
func AdminGetVacancyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	v, err := repository.GetVacancyByID(r.Context(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if v == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// AdminCreateVacancy — создаёт вакансию.
func AdminCreateVacancy(w http.ResponseWriter, r *http.Request) {
	var p models.VacancyPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}
	v := &models.Vacancy{
		Title:        p.Title,
		Description:  p.Description,
		Requirements: p.Requirements,
		Conditions:   p.Conditions,
		ContactEmail: p.ContactEmail,
	}
	id, err := repository.CreateVacancy(r.Context(), v)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// AdminUpdateVacancy — обновляет вакансию по id.
func AdminUpdateVacancy(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	var p models.VacancyPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}
	v := &models.Vacancy{
		ID:           id,
		Title:        p.Title,
		Description:  p.Description,
		Requirements: p.Requirements,
		Conditions:   p.Conditions,
		ContactEmail: p.ContactEmail,
	}
	if err := repository.UpdateVacancy(r.Context(), v); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminDeleteVacancy — удаляет вакансию.
func AdminDeleteVacancy(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteVacancy(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
