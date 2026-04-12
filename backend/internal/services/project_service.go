package services

import (
	"database/sql"
	"errors"
	"time"

	"taskflow/internal/models"

	"github.com/google/uuid"
)

type ProjectService struct {
	DB *sql.DB
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{DB: db}
}

// CREATE PROJECT
func (s *ProjectService) CreateProject(name, description, ownerID string) (*models.Project, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}

	id := uuid.New().String()
	now := time.Now()

	_, err := s.DB.Exec(`
		INSERT INTO projects (id, name, description, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, name, description, ownerID, now)

	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:          id,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   now,
	}, nil
}

// GET ALL PROJECTS (USER)
func (s *ProjectService) GetProjectsByUser(userID string) ([]models.Project, error) {
	rows, err := s.DB.Query(`
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.OwnerID,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// GET PROJECT BY ID
func (s *ProjectService) GetProjectByID(projectID, userID string) (*models.Project, error) {
	var p models.Project

	err := s.DB.QueryRow(`
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE id = $1 AND owner_id = $2
	`, projectID, userID).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.OwnerID,
		&p.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return &p, nil
}

// UPDATE PROJECT
func (s *ProjectService) UpdateProject(projectID, name, description, userID string) error {
	if name == "" {
		return errors.New("project name is required")
	}

	res, err := s.DB.Exec(`
		UPDATE projects
		SET name = $1,
		    description = $2
		WHERE id = $3 AND owner_id = $4
	`, name, description, projectID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("project not found or unauthorized")
	}

	return nil
}

// DELETE PROJECT
func (s *ProjectService) DeleteProject(projectID, userID string) error {
	res, err := s.DB.Exec(`
		DELETE FROM projects
		WHERE id = $1 AND owner_id = $2
	`, projectID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("project not found or unauthorized")
	}

	return nil
}
