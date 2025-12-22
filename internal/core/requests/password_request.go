package requests

import "gophkeep/internal/core/models"

type PasswordRequestObject struct {
	Login    string        `json:"login"`
	Password string        `json:"password"`
	UserID   models.UserID `json:"user_id"`
	Comment  string        `json:"comment"`
}
