package dto

// RegisterRequest represents payload for user registration
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=150"`
	FirstName string `json:"first_name" binding:"required,min=3,max=150"`
	LastName  string `json:"last_name" binding:"required,min=3,max=150"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents payload for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// LoginResponse represents successful login with token
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// AuthResponse returns user + token (kept for compatibility)
type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
