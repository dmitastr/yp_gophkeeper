package models

import "errors"

type Password struct {
	ID       int    `json:"secret_id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Comment  string `json:"comment,omitempty"`
}

func (p *Password) Validate() error {
	if p.Login == "" || p.Password == "" {
		return errors.New("login or password is empty")
	}
	return nil
}
