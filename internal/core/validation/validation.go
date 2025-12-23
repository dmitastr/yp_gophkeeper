package validation

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
)

type SecretHandler func([]byte) ([]byte, error)

type IValidator interface {
	Validate(secretContent string, secretType models.SecretType, isEncoded bool) ([]byte, error)
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

func (v *Validator) parseTextSecret(content []byte) ([]byte, error) {
	if utf8.Valid(content) {
		return content, nil
	}
	return nil, ErrorInvalidSecretContent
}

func (v *Validator) parsePasswordSecret(content []byte) ([]byte, error) {
	var r models.Password
	if err := json.Unmarshal(content, &r); err != nil {
		return nil, fmt.Errorf("error unmarshalling password secret: %w", err)
	}
	return content, r.Validate()
}

func (v *Validator) parseBankCardSecret(content []byte) ([]byte, error) {
	v.cfg.Logger().Info("Validating bank card info", zap.String("message", string(content)))
	var r models.BankCard
	if err := json.Unmarshal(content, &r); err != nil {
		return nil, fmt.Errorf("error unmarshalling bank card secret: %w", err)
	}
	return content, r.Validate()
}

func (v *Validator) parseBinarySecret(content []byte) ([]byte, error) {
	if len(content) == 0 || content == nil {
		return nil, ErrorEmptySecretContent
	}
	return content, nil
}

func (v *Validator) decode(message string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("invalid binary secret")
	}
	return content, nil
}

func (v *Validator) Validate(secretContent string, secretType models.SecretType, isEncoded bool) ([]byte, error) {
	var content []byte
	var err error

	if isEncoded {
		content, err = v.decode(secretContent)
		if err != nil {
			return nil, err
		}
	} else {
		content = []byte(secretContent)
	}

	handler, ok := v.secretRegistry[secretType]
	if !ok {

		return nil, ErrorSecretTypeNotFound
	}
	secret, err := handler(content)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
