package handlers

import (
	"encoding/json"
	"net/http"
	"x86trade_backend/internal/repository"
)

func GetPaymentMethodsHandler(w http.ResponseWriter, r *http.Request) {
	methods, err := repository.GetPaymentMethods(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}
