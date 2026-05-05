package models

type PasswordRequestObject struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
