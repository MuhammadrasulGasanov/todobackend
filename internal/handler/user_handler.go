package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"id":       user.ID,
		"username": user.Username,
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
