package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"x86trade_backend/models"
	"x86trade_backend/repository"

	"github.com/go-chi/chi/v5"
)

func GetProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		http.Error(w, "bad request: invalid product id", http.StatusBadRequest)
		return
	}

	log.Printf("Getting details for product ID: %d", productID)

	productDetail, err := repository.GetProductDetails(r.Context(), productID)
	if err != nil {
		log.Printf("Error getting product details: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if productDetail == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productDetail)
}

func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID int    `json:"product_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}

	if req.ProductID <= 0 {
		http.Error(w, "bad request: invalid product id", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "bad request: rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	if len(req.Comment) < 10 {
		http.Error(w, "bad request: comment must be at least 10 characters", http.StatusBadRequest)
		return
	}

	review := &models.Review{
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := repository.CreateReview(r.Context(), review); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Review created successfully",
	})
}

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

// AdminGetProductCharacteristics — переиспользует GetProductCharacteristics (уже есть)
func AdminGetProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, _ := strconv.Atoi(productIDStr)
	if productID <= 0 {
		http.Error(w, "bad request: product id", http.StatusBadRequest)
		return
	}
	chars, err := repository.GetProductCharacteristics(r.Context(), productID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(chars)
}

// Create one characteristic (for product)
func AdminCreateProductCharacteristic(w http.ResponseWriter, r *http.Request) {
	var p models.ProductCharacteristicInput
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.ProductID <= 0 || p.CharacteristicTypeID <= 0 || p.Value == "" {
		http.Error(w, "product_id, characteristic_type_id and value required", http.StatusBadRequest)
		return
	}
	id, err := repository.CreateProductCharacteristic(r.Context(), &p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func AdminUpdateProductCharacteristic(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	var p models.ProductCharacteristicInput
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p.ID = id
	if p.ID <= 0 || p.CharacteristicTypeID <= 0 || p.Value == "" {
		http.Error(w, "id, characteristic_type_id and value required", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateProductCharacteristic(r.Context(), &p); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AdminDeleteProductCharacteristic(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		http.Error(w, "bad request: id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteProductCharacteristic(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ReplaceAllProductCharacteristics — принимает массив и заменяет все характеристики товара (transactional)
func AdminReplaceProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, _ := strconv.Atoi(productIDStr)
	if productID <= 0 {
		http.Error(w, "bad request: product id", http.StatusBadRequest)
		return
	}
	var inputs []models.ProductCharacteristicInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// Валидация: каждый элемент должен содержать characteristic_type_id и value
	for _, it := range inputs {
		if it.CharacteristicTypeID <= 0 || it.Value == "" {
			http.Error(w, "each characteristic must include characteristic_type_id and non-empty value", http.StatusBadRequest)
			return
		}
		it.ProductID = productID
	}
	if err := repository.ReplaceProductCharacteristics(r.Context(), productID, inputs); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
