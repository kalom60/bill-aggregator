package models

import "time"

type Provider struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name" validate:"required"`
	API_URL             string    `json:"api_url" validate:"required"`
	Authentication_Type string    `json:"authentication_type" validate:"required"`
	API_key             string    `json:"api_key" validate:"required"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
