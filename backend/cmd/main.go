package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"taskflow/internal/db"
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"
	"taskflow/internal/services"
)

func main() {
	// Database connection
	db.Connect()

	// Router setup
	r := chi.NewRouter()

	//  HANDLE PREFLIGHT (VERY IMPORTANT)
	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//  CORS (VERY IMPORTANT for React)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Post("/auth/register", handlers.RegisterHandler)
	r.Post("/auth/login", handlers.LoginHandler)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Services and handlers
	projectService := services.NewProjectService(db.DB)
	projectHandler := handlers.NewProjectHandler(projectService)

	taskService := services.NewTaskService(db.DB)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Protected routes
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)

		protected.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You are authorized"))
		})

		// Project routes
		protected.Post("/projects", projectHandler.CreateProject)
		protected.Get("/projects", projectHandler.GetProjects)
		protected.Get("/projects/{id}", projectHandler.GetProjectByID)
		protected.Patch("/projects/{id}", projectHandler.UpdateProject)
		protected.Delete("/projects/{id}", projectHandler.DeleteProject)

		// Task routes
		protected.Post("/projects/{id}/tasks", taskHandler.CreateTask)
		protected.Get("/projects/{id}/tasks", taskHandler.GetTasks)
		protected.Patch("/tasks/{id}", taskHandler.UpdateTask)
		protected.Delete("/tasks/{id}", taskHandler.DeleteTask)
	})

	// Server start
	log.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
