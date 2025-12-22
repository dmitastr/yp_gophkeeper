package requests

import (
	"encoding/json"

	"gophkeep/internal/core/models"
)

type SecretRequest struct {
	ID         int               `json:"id,omitempty"`
	Body       json.RawMessage   `json:"body"`
	BodyString string            `json:"body_string"`
	SecretType models.SecretType `json:"secret_type"`
	Comment    string            `json:"comment"`
}

type PasswordRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
