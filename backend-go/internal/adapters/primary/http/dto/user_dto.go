package dto

// RegisterRequest represents payload for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=150"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents payload for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse returns user + token
type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
