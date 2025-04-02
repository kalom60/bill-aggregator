package models

import "time"

type Account struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id" validate:"required"`
	ProviderID          string    `json:"provider_id" validate:"required"`
	AccountIdentifier   string    `json:"account_identifier" validate:"required"`
	EncryptedCredential string    `json:"encrypted_credential" validate:"required"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
