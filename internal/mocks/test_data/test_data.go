package test_data

import (
	"encoding/json"

	"gophkeep/internal/core/models"
)

var PasswordTestData = &models.Password{
	Login:    "login",
	Password: "password",
	Comment:  "comment",
}

var BankcardTestData, _ = json.Marshal(&models.BankCard{
	Number:         5120350100064537,
	ExpiryData:     "03/30",
	CardholderName: "Ivan",
	CVC:            "000",
	Comment:        "comment",
})

var TextTestData = []byte("text secret")

var BinaryTestData = []byte("binary secret")

type TestData map[models.SecretType][]byte

func GetTestData() TestData {
	t := make(map[models.SecretType][]byte)

	passwordData, _ := json.Marshal(PasswordTestData)
	t[models.PASSWORD] = passwordData
	bankcardData, _ := json.Marshal(BankcardTestData)
	t[models.BANK_CARD] = bankcardData

	t[models.TEXT] = TextTestData
	t[models.BINARY] = BinaryTestData

	return t
}
