package docs

import "time"

// ErrorResponse represents a standard API error response
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// RegisterRequest represents the payload for user registration
type RegisterRequest struct {
	Username string `json:"username" example:"johndoe" validate:"required,min=3,max=150"`
	Email    string `json:"email" example:"john@example.com" validate:"required,email"`
	Password string `json:"password" example:"SecurePass123" validate:"required,min=8"`
}

// LoginRequest represents the payload for user login
type LoginRequest struct {
	Username string `json:"username" example:"johndoe" validate:"required"`
	Password string `json:"password" example:"SecurePass123" validate:"required"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	FirstName string    `json:"first_name,omitempty" example:"John"`
	LastName  string    `json:"last_name,omitempty" example:"Doe"`
	Role      string    `json:"role" example:"cliente"`
	CreatedAt time.Time `json:"created_at" example:"2025-11-07T10:00:00Z"`
	LastLogin *time.Time `json:"last_login" example:"2025-11-07T10:00:00Z"`
	IsActive  bool      `json:"is_active" example:"true"`
}

// AuthResponse represents successful authentication responses
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJ..."`
}

// @Summary Register new user
// @Description Creates a new user account and returns JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "User registration data"
// @Success 201 {object} AuthResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Router /api/users/register [post]

// @Summary Login user
// @Description Authenticates user and returns JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param data body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse "Login successful"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /api/users/login [post]

// @Summary Get user profile
// @Description Returns the authenticated user's profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse "User profile"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /api/users/profile [get]
