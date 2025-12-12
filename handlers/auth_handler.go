package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"x86trade_backend/models"
	"x86trade_backend/repository"
	"x86trade_backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// helpers to read envs
func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return def
}

// Register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.Email == "" || payload.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}
	// check existing
	if existing, _ := repository.GetUserByEmail(r.Context(), payload.Email); existing != nil {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	}
	// hash pwd
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	user := &models.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}
	id, err := repository.CreateUser(r.Context(), user, string(hashed))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, err := repository.GetUserByEmail(r.Context(), payload.Email)
	if err != nil || u == nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(payload.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// generate access token
	accessMinutes := getEnvInt("JWT_ACCESS_MINUTES", 15)
	accessToken, err := utils.GenerateAccessToken(u.ID, accessMinutes)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// generate refresh token (random string) and save
	refreshDays := getEnvInt("JWT_REFRESH_DAYS", 7)
	// random 32 bytes -> hex
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	refreshToken := hex.EncodeToString(b)
	expiresAt := time.Now().Add(time.Duration(refreshDays) * 24 * time.Hour)
	if err := repository.SaveRefreshToken(r.Context(), u.ID, refreshToken, expiresAt); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// response
	resp := map[string]interface{}{
		"access_token":   accessToken,
		"token_type":     "bearer",
		"expires_in_min": accessMinutes,
		"refresh_token":  refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Refresh
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.RefreshToken == "" {
		http.Error(w, "refresh_token required", http.StatusBadRequest)
		return
	}
	userID, expiresAt, err := repository.GetRefreshToken(r.Context(), payload.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	if time.Now().After(expiresAt) {
		// token expired — delete and ask to login again
		_ = repository.DeleteRefreshToken(r.Context(), payload.RefreshToken)
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}

	accessMinutes := getEnvInt("JWT_ACCESS_MINUTES", 15)
	accessToken, err := utils.GenerateAccessToken(userID, accessMinutes)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":   accessToken,
		"token_type":     "bearer",
		"expires_in_min": accessMinutes,
	})
}

// Logout: delete refresh token
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if payload.RefreshToken == "" {
		http.Error(w, "refresh_token required", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteRefreshToken(r.Context(), payload.RefreshToken); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
