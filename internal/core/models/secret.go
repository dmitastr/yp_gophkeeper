package models

import (
	"errors"
	"time"
)

type SecretType string

const (
	PASSWORD  SecretType = "password"
	TEXT      SecretType = "text"
	BANK_CARD SecretType = "bank_card"
	BINARY    SecretType = "binary"
)

type SecretTypeValue struct {
	Value *SecretType
}

func (st *SecretTypeValue) String() string {
	if st.Value == nil {
		return ""
	}
	return string(*st.Value)
}

func (st *SecretTypeValue) Set(value string) error {
	switch tp := SecretType(value); tp {
	case PASSWORD, TEXT, BANK_CARD, BINARY:
		*st.Value = SecretType(value)
	default:
		return errors.New("invalid secret type")
	}
	return nil
}

func (st *SecretTypeValue) Type() string {
	return "secretType"
}

type Secret struct {
	Type      SecretType `json:"type" db:"secret_type"`
	ID        int        `json:"id" db:"secret_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Content   []byte     `json:"content,omitempty" db:"secret"`
	Comment   string     `json:"comment,omitempty" db:"comment"`
}

type SecretInfo struct {
	Type      SecretType `json:"type" db:"secret_type"`
	ID        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
