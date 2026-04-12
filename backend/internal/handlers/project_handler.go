package handlers

import (
	"encoding/json"
	"net/http"

	"taskflow/internal/services"

	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	Service *services.ProjectService
}

func NewProjectHandler(s *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{Service: s}
}

// HELPER RESPONSE METHODS

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}

// CREATE PROJECT

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.Service.CreateProject(req.Name, req.Description, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, project)
}

// GET ALL PROJECTS

func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	projects, err := h.Service.GetProjectsByUser(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}

	respondJSON(w, http.StatusOK, projects)
}

// GET PROJECT BY ID

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	projectID := chi.URLParam(r, "id")

	project, err := h.Service.GetProjectByID(projectID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, project)
}

// UPDATE PROJECT

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	projectID := chi.URLParam(r, "id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.Service.UpdateProject(projectID, req.Name, req.Description, userID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "project updated successfully",
	})
}

// DELETE PROJECT

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	projectID := chi.URLParam(r, "id")

	err := h.Service.DeleteProject(projectID, userID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "project deleted successfully",
	})
}
