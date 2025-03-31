package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8"`
	FirstName string    `json:"first_name" validate:"required,matches=^[A-Za-z ]+$"`
	LastName  string    `json:"last_name" validate:"required,matches=^[A-Za-z ]+$"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
