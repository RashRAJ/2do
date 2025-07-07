package models

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
type CreateUserRequest struct {
	Username string `json:"username" validate:"required, min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required, min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
