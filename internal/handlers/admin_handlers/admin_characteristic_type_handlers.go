package admin_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"x86trade_backend/internal/models"
	"x86trade_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func AdminGetCharacteristicTypes(w http.ResponseWriter, r *http.Request) {
	types, err := repository.GetAllCharacteristicTypes(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(types)
}

func AdminCreateCharacteristicType(w http.ResponseWriter, r *http.Request) {
	var p models.CharacteristicType
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateCharacteristicType(r.Context(), &p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func AdminUpdateCharacteristicType(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	var p models.CharacteristicType
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p.ID = id
	if p.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateCharacteristicType(r.Context(), &p); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AdminDeleteCharacteristicType(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if err := repository.DeleteCharacteristicType(r.Context(), id); err != nil {
		// если FK -> зависимые записи, БД вернёт ошибку; её можно перехватить и вернуть 409 Conflict
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
