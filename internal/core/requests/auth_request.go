package requests

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IsNewUser bool   `json:"-"`
}
