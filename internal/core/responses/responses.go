package responses

import "gophkeep/internal/core/models"

type SecretsListResponseObject struct {
	Secrets []models.SecretInfo `json:"secrets"`
}

type SecretResponseObject struct {
	Secret *models.Secret `json:"secret"`
}
