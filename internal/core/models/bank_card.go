package models

import (
	"errors"

	"github.com/theplant/luhn"
)

type BankCard struct {
	ID             int    `json:"secret_id,omitempty" db:"secret_id"`
	Number         int    `json:"number"`
	ExpiryData     string `json:"expiry_data"`
	CardholderName string `json:"cardholder_name"`
	CVC            string `json:"cvc,omitempty"`
	Comment        string `json:"comment,omitempty"`
}

func (b *BankCard) IsValid() bool {
	return luhn.Valid(b.Number)
}

func (b *BankCard) Validate() error {
	if !b.IsValid() {
		return errors.New("card number is invalid")
	}
	return nil
}
