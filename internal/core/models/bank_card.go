package models

import "github.com/theplant/luhn"

type BankCard struct {
	ID             int    `json:"secret_id" db:"secret_id"`
	Number         int    `json:"number"`
	ExpiryData     string `json:"expiry_data"`
	CardholderName string `json:"cardholder_name"`
	CVC            string `json:"cvc,omitempty"`
}

func (b *BankCard) IsValid() bool {
	return luhn.Valid(b.Number)
}
