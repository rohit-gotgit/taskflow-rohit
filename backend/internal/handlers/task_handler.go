package handlers

import (
	"encoding/json"
	"net/http"

	"taskflow/internal/services"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(s *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
}

// CREATE TASK
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
		AssigneeID  string `json:"assignee_id"`
	}

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, 400, "invalid request body")
		return
	}

	task, errs, err := h.Service.CreateTask(
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		projectID,
		req.AssigneeID,
	)

	// server error
	if err != nil {
		respondError(w, 500, "server error")
		return
	}

	// validation error
	if errs != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
		return
	}

	respondJSON(w, 201, task)
}

// GET TASKS
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")

	tasks, err := h.Service.GetTasks(projectID, status, assignee)
	if err != nil {
		respondError(w, 500, "failed to fetch tasks")
		return
	}

	respondJSON(w, 200, tasks)
}

// UPDATE TASK
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	// safer context extraction 🔥
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondError(w, 401, "unauthorized")
		return
	}

	// ownership check
	owned, err := h.Service.IsTaskOwnedByUser(taskID, userID)
	if err != nil {
		respondError(w, 500, "server error")
		return
	}
	if !owned {
		respondError(w, 403, "forbidden")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, 400, "invalid request body")
		return
	}

	errs, err := h.Service.UpdateTask(taskID, updates)

	// server error
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}

	// validation error
	if errs != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
		return
	}

	respondJSON(w, 200, map[string]string{
		"message": "task updated",
	})
}

// DELETE TASK
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	// safer context extraction
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondError(w, 401, "unauthorized")
		return
	}

	// ownership check
	owned, err := h.Service.IsTaskOwnedByUser(taskID, userID)
	if err != nil {
		respondError(w, 500, "server error")
		return
	}
	if !owned {
		respondError(w, 403, "forbidden")
		return
	}

	err = h.Service.DeleteTask(taskID)
	if err != nil {
		respondError(w, 500, "failed to delete task")
		return
	}

	respondJSON(w, 200, map[string]string{
		"message": "task deleted",
	})
}
