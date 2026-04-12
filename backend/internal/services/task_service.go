package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"taskflow/internal/models"

	"github.com/google/uuid"
)

type TaskService struct {
	DB *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{DB: db}
}

// VALIDATION MAPS (GLOBAL)
var validStatus = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"done":        true,
}

var validPriority = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

// VALIDATION FUNCTION (CENTRALIZED)
func validateTaskInput(title, status, priority string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "required"
	}

	if status != "" && !validStatus[status] {
		errors["status"] = "invalid value"
	}

	if priority != "" && !validPriority[priority] {
		errors["priority"] = "invalid value"
	}

	return errors
}

// CREATE TASK
func (s *TaskService) CreateTask(title, description, status, priority, projectID, assigneeID string) (*models.Task, map[string]string, error) {

	// validation
	errs := validateTaskInput(title, status, priority)
	if len(errs) > 0 {
		return nil, errs, nil
	}

	id := uuid.New().String()
	now := time.Now()

	var assignee interface{}
	if assigneeID == "" {
		assignee = nil
	} else {
		assignee = assigneeID
	}

	_, err := s.DB.Exec(`
	INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, id, title, description, status, priority, projectID, assignee, now, now)

	if err != nil {
		return nil, nil, err
	}

	var assigneePtr *string
	if assigneeID != "" {
		assigneePtr = &assigneeID
	}

	return &models.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		ProjectID:   projectID,
		AssigneeID:  assigneePtr,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil, nil
}

// GET TASKS
func (s *TaskService) GetTasks(projectID, status, assignee string) ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, project_id, assignee_id, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	args := []interface{}{projectID}
	i := 2

	if status != "" {
		query += " AND status = $" + fmt.Sprint(i)
		args = append(args, status)
		i++
	}

	if assignee != "" {
		query += " AND assignee_id = $" + fmt.Sprint(i)
		args = append(args, assignee)
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.ProjectID,
			&t.AssigneeID,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// SAFE UPDATE TASK (WHITELIST + VALIDATION)
func (s *TaskService) UpdateTask(taskID string, updates map[string]interface{}) (map[string]string, error) {

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	// whitelist fields
	allowed := map[string]bool{
		"title":       true,
		"description": true,
		"status":      true,
		"priority":    true,
	}

	setClauses := []string{}
	args := []interface{}{}
	i := 1

	validationErrors := make(map[string]string)

	for key, value := range updates {

		if !allowed[key] {
			continue
		}

		// validation for specific fields
		if key == "status" {
			val := value.(string)
			if !validStatus[val] {
				validationErrors["status"] = "invalid value"
			}
		}

		if key == "priority" {
			val := value.(string)
			if !validPriority[val] {
				validationErrors["priority"] = "invalid value"
			}
		}

		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", key, i))
		args = append(args, value)
		i++
	}

	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no valid fields to update")
	}

	query := fmt.Sprintf(
		"UPDATE tasks SET %s, updated_at=$%d WHERE id=$%d",
		strings.Join(setClauses, ", "),
		i,
		i+1,
	)

	args = append(args, time.Now(), taskID)

	_, err := s.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DELETE TASK
func (s *TaskService) DeleteTask(taskID string) error {
	_, err := s.DB.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	return err
}

// OWNERSHIP CHECK
func (s *TaskService) IsTaskOwnedByUser(taskID, userID string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tasks t
			JOIN projects p ON t.project_id = p.id
			WHERE t.id = $1 AND p.owner_id = $2
		)
	`

	err := s.DB.QueryRow(query, taskID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
