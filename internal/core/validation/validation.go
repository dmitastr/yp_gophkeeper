package validation

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"

	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
)

type SecretHandler func(string) ([]byte, error)

type IValidator interface {
	Validate(secretContent string, secretType models.SecretType) ([]byte, error)
}

type Validator struct {
	secretRegistry map[models.SecretType]SecretHandler
	cfg            config.ConfigProvider
}

func NewValidator(cfg config.ConfigProvider) IValidator {
	v := &Validator{cfg: cfg}
	v.secretRegistry = map[models.SecretType]SecretHandler{
		models.PASSWORD:  v.parsePasswordSecret,
		models.BANK_CARD: v.parseBankCardSecret,
		models.TEXT:      v.parseTextSecret,
		models.BINARY:    v.parseBinarySecret,
	}
	return v
}

func (v *Validator) Validate(secretContent string, secretType models.SecretType) ([]byte, error) {
	handler, ok := v.secretRegistry[secretType]
	if !ok {
		return nil, errors.New("secretType not exist")
	}
	secret, err := handler(secretContent)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (v *Validator) parseTextSecret(message string) ([]byte, error) {
	content := []byte(message)

	if utf8.Valid(content) {
		return content, nil
	}
	return nil, errors.New("invalid message")
}

func (v *Validator) parsePasswordSecret(message string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding message: %w", err)
	}

	var r models.Password
	if err := json.Unmarshal(content, &r); err != nil {
		return nil, fmt.Errorf("error unmarshalling password secret: %w", err)
	}
	return content, r.Validate()
}

func (v *Validator) parseBankCardSecret(message string) ([]byte, error) {
	v.cfg.Logger().Info("Validating bank card info", zap.String("message", message))
	content, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding message: %w", err)
	}
	var r models.BankCard
	if err := json.Unmarshal(content, &r); err != nil {
		return nil, fmt.Errorf("error unmarshalling bank card secret: %w", err)
	}
	return content, r.Validate()
}

func (v *Validator) parseBinarySecret(message string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("invalid binary secret")
	}
	if len(content) == 0 || content == nil {
		return nil, errors.New("binary data is empty")
	}
	return content, nil
}
