package models

import "time"

type SecretType string

const (
	PASSWORD  SecretType = "password"
	TEXT      SecretType = "text"
	BANK_CARD SecretType = "bank_card"
	BINARY    SecretType = "binary"
)

type Secret struct {
	Type      SecretType `json:"type" db:"secret_type"`
	ID        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Content   []byte     `json:"content,omitempty" db:"secret"`
	Comment   string     `json:"comment,omitempty" db:"comment"`
}

type SecretInfo struct {
	Type      SecretType `json:"type" db:"secret_type"`
	ID        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
