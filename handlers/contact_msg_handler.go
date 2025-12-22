package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"x86trade_backend/models"
	"x86trade_backend/repository"
)

func CreateContactMessageHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FullName    string `json:"full_name"`
		ContactInfo string `json:"contact_info"`
		Message     string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if payload.FullName == "" || payload.ContactInfo == "" || payload.Message == "" {
		http.Error(w, "bad request: all fields are required", http.StatusBadRequest)
		return
	}

	// Создаем сообщение
	msg := &models.ContactMessage{
		FullName:    payload.FullName,
		ContactInfo: payload.ContactInfo,
		Message:     payload.Message,
		CreatedAt:   time.Now(),
		IsProcessed: false,
	}

	id, err := repository.CreateContactMessage(r.Context(), msg)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Ваше сообщение успешно отправлено",
	})
}
