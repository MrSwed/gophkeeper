package model

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Meta     any    `json:"meta"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
