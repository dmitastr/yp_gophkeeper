package requests

import (
	"encoding/json"

	"gophkeep/internal/core/models"
)

type SecretRequest struct {
	Body       json.RawMessage   `json:"body"`
	SecretType models.SecretType `json:"secret_type"`
	Comment    string            `json:"comment"`
}

type PasswordRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
