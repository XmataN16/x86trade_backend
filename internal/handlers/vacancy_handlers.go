package handlers

import (
	"encoding/json"
	"net/http"
	"x86trade_backend/internal/repository"
)

func GetVacanciesHandler(w http.ResponseWriter, r *http.Request) {
	vacancies, err := repository.GetVacancies(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vacancies)
}
