package handlers

import (
	"encoding/json"
	"net/http"

	"x86trade_backend/internal/repository"
)

func GetDeliveryMethodsHandler(w http.ResponseWriter, r *http.Request) {
	methods, err := repository.GetDeliveryMethods(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// Convert sql.Null* to simple JSON fields (we can reuse structs but simplest: encode directly)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}
