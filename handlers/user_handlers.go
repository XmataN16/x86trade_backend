package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"x86trade_backend/models"
	"x86trade_backend/repository"
)

// AdminGetUsers — возвращает всех пользователей (без password_hash).
func AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// AdminCreateUser — создаёт пользователя (принимает пароль в теле).
// JSON: { "email": "...", "password": "...", "first_name": "...", "last_name": "...", "phone": "...", "is_admin": true/false }
func AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		IsAdmin   bool   `json:"is_admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.Email == "" || payload.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}
	// hash password
	h, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	u := &models.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Phone:     payload.Phone,
		IsAdmin:   payload.IsAdmin,
	}
	id, err := repository.CreateUser(r.Context(), u, string(h))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// AdminUpdateUser — обновляет данные пользователя (не меняет пароль).
// JSON: { "email": "...", "first_name": "...", "last_name": "...", "phone": "...", "is_admin": true/false }
func AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	var payload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		IsAdmin   bool   `json:"is_admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u := &models.User{
		ID:        id,
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Phone:     payload.Phone,
		IsAdmin:   payload.IsAdmin,
	}
	if err := repository.UpdateUser(r.Context(), u); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminDeleteUser — удаляет пользователя по id.
func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteUser(r.Context(), id); err != nil {
		log.Printf("AdminDeleteUser error id=%d: %v", id, err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			// FK violation
			http.Error(w, "conflict: cannot delete, resource is referenced by other records", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AdminGetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	u, err := repository.GetUserByID(r.Context(), id)
	if err != nil {
		log.Printf("AdminGetUserByID error id=%d: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if u == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// AdminUpdateUserPassword обновляет пароль пользователя
func AdminUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var payload struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	if len(payload.Password) < 6 {
		http.Error(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := repository.UpdateUserPassword(r.Context(), userID, string(hashed)); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
