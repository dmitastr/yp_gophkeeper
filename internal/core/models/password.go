package models

type Password struct {
	ID       int    `json:"secret_id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Comment  string `json:"comment,omitempty"`
}

func (p *Password) Validate() error {
	if p.Login == "" || p.Password == "" {
		return ErrorEmptyAuthData
	}
	return nil
}
