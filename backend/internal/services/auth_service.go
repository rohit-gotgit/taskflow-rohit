package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"taskflow/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("taskflow_secret") // should match middleware

// Register a new user
func Register(name, email, password string) error {
	// basic validation
	if name == "" || email == "" || password == "" {
		return errors.New("all fields are required")
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// insert user
	_, err = db.DB.Exec(
		"INSERT INTO users (name, email, password) VALUES ($1,$2,$3)",
		name, email, string(hashed),
	)

	if err != nil {
		// handle duplicate email
		if strings.Contains(err.Error(), "users_email_key") {
			return errors.New("email already exists")
		}
		return err
	}

	return nil
}

// Login user and return JWT token
func Login(email, password string) (string, error) {
	// basic validation
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	var id, hashed string

	// fetch user
	err := db.DB.QueryRow(
		"SELECT id, password FROM users WHERE email=$1",
		email,
	).Scan(&id, &hashed)

	if err == sql.ErrNoRows {
		return "", errors.New("invalid credentials")
	}
	if err != nil {
		return "", err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	// sign token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
