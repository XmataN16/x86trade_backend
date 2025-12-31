package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"x86trade_backend/models"
	"x86trade_backend/repository"

	"github.com/go-chi/chi/v5"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cats, err := repository.GetCategories(ctx)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	cat, err := repository.GetCategoryByID(context.Background(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if cat == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateCategory(context.Background(), &c)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}
	if c.ID == 0 {
		http.Error(w, "bad request: id required", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateCategory(context.Background(), &c); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteCategory(context.Background(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminGetCategories — список категорий.
func AdminGetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := repository.GetAllCategories(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cats)
}

// AdminGetCategory возвращает категорию по ID
func AdminGetCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		http.Error(w, "bad request: invalid category id", http.StatusBadRequest)
		return
	}

	c, err := repository.GetCategoryByID(r.Context(), categoryID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if c == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// AdminCreateCategory — создать категорию.
// JSON: { "name":"...","slug":"..." }
func AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload models.Category
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateCategory(r.Context(), &payload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// AdminUpdateCategory — обновление категории.
// JSON: { "name":"...","slug":"..." }
func AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	var payload models.Category
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	payload.ID = id
	if payload.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateCategory(r.Context(), &payload); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminDeleteCategory — удалить категорию по id.
func AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteCategory(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
