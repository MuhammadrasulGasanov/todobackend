package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/middleware"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/models"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/service"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	Service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		CategoryID  *int    `json:"category_id"`
		DueDate     *string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if input.DueDate != nil {
		parsed, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
		dueDate = &parsed
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		http.Error(w, "due date cannot be in the past", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		CategoryID:  input.CategoryID,
		DueDate:     dueDate,
		Completed:   false,
	}

	if err := h.Service.CreateTask(r.Context(), task); err != nil {
		http.Error(w, "could not create task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	categoryIDStr := r.URL.Query().Get("category_id")
	var categoryID *int
	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "invalid category_id", http.StatusBadRequest)
			return
		}
		categoryID = &id
	}

	tasks, err := h.Service.GetTasksByUser(r.Context(), userID, categoryID)
	if err != nil {
		http.Error(w, "could not get tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}
	var input struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		CategoryID  *int    `json:"category_id"`
		DueDate     *string `json:"due_date"`
		Completed   bool    `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if input.DueDate != nil {
		parsed, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
		dueDate = &parsed
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		http.Error(w, "due date cannot be in the past", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		CategoryID:  input.CategoryID,
		DueDate:     dueDate,
		Completed:   input.Completed,
	}

	if err := h.Service.UpdateTask(r.Context(), task); err != nil {
		http.Error(w, "could not update task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) MarkTaskCompletion(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}
	var input struct {
		Completed bool `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.MarkTaskCompletion(r.Context(), taskID, userID, input.Completed); err != nil {
		http.Error(w, "could not update task completion", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.Service.GetTaskByID(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteTask(r.Context(), taskID, userID); err != nil {
		http.Error(w, "could not delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
