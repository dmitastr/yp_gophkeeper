package models

type Password struct {
	ID       int
	Login    string `json:"login"`
	Password string `json:"password"`
	Source   string `json:"source"`
}
