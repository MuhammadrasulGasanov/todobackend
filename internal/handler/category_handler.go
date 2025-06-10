package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/middleware"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/models"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/service"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	Service service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: *s}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	category := &models.Category{
		Name:   input.Name,
		UserID: userID,
	}

	if err := h.Service.CreateCategory(r.Context(), category); err != nil {
		http.Error(w, "could not create category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteCategory(r.Context(), categoryID, userID); err != nil {
		http.Error(w, "could not delete category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	categories, err := h.Service.GetCategoriesByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "could not get categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)

}

func (h *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.Service.GetCategoryById(r.Context(), categoryID, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}
