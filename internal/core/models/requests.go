package models

type PasswordRequestObject struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   UserID `json:"user_id"`
	Comment  string `json:"comment"`
}
