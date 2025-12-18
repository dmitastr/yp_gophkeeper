package requests

import "gophkeep/internal/core/models"

type SecretRequest struct {
	Body       []byte            `json:"body"`
	SecretType models.SecretType `json:"secret_type"`
	Comment    string            `json:"comment"`
}

type PasswordRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
