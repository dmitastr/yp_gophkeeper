package models

type Password struct {
	ID       int    `json:"secret_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Source   string `json:"source"`
}
