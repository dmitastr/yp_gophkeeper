package validation

import (
	"encoding/json"
	"errors"
	"unicode/utf8"

	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
)

type SecretHandler func(json.RawMessage) (any, error)

type IValidator interface {
	Validate(secretRequest *requests.SecretRequest) (any, error)
}

type Validator struct {
	secretRegistry map[models.SecretType]SecretHandler
}

func NewValidator() IValidator {
	v := &Validator{}
	v.secretRegistry = map[models.SecretType]SecretHandler{
		models.PASSWORD:  v.parsePasswordSecret,
		models.BANK_CARD: v.parseBankCardSecret,
		models.TEXT:      v.parseTextSecret,
		models.BINARY:    v.parseBinarySecret,
	}
	return v
}

func (v *Validator) Validate(secretRequest *requests.SecretRequest) (any, error) {
	handler, ok := v.secretRegistry[secretRequest.SecretType]
	if !ok {
		return nil, errors.New("secretType not exist")
	}
	secret, err := handler(secretRequest.Body)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (v *Validator) parseTextSecret(message json.RawMessage) (any, error) {
	if utf8.Valid(message) {
		return message, nil
	}
	return nil, errors.New("invalid message")
}

func (v *Validator) parsePasswordSecret(message json.RawMessage) (any, error) {
	var r models.Password
	if err := json.Unmarshal(message, &r); err != nil {
		return nil, err
	}
	return &r, r.Validate()
}

func (v *Validator) parseBankCardSecret(message json.RawMessage) (any, error) {
	var r models.BankCard
	if err := json.Unmarshal(message, &r); err != nil {
		return nil, err
	}
	return &r, r.Validate()
}

func (v *Validator) parseBinarySecret(message json.RawMessage) (any, error) {
	if len(message) == 0 || message == nil {
		return nil, errors.New("binary data is empty")
	}
	return message, nil
}
