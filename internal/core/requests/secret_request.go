package requests

import (
	"gophkeep/internal/core/models"
)

type SecretRequest struct {
	ID         *int              `json:"id,omitempty"`
	BodyString string            `json:"body_string"`
	SecretType models.SecretType `json:"secret_type"`
	Comment    string            `json:"comment"`
	IsEncoded  bool              `json:"is_encoded"`
}

type PasswordRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
