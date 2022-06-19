package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Password_again string `json:"password_again"`
}
