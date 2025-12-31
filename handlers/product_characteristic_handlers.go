package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"x86trade_backend/repository"
)

func AdminListProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	offset := (page - 1) * limit

	list, err := repository.GetAllProductCharacteristicsWithPagination(r.Context(), limit, offset)
	if err != nil {
		// логируем в stdout/stderr (поможет диагностировать 500)
		// log.Printf("AdminListProductCharacteristics error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	total, err := repository.CountProductCharacteristics(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"data":  list,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
