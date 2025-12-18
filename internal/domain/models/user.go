package models

type UserID int

type User struct {
	Username string `json:"username"`
	ID       UserID `json:"user_id"`
	Password string `json:"password"`
	Hash     string `json:"hash"`
}
