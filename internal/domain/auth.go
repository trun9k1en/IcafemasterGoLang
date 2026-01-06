package domain

import "context"

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	User         *UserInfo `json:"user"`
}

// UserInfo represents user info in token response
type UserInfo struct {
	ID          string       `json:"id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	FullName    string       `json:"full_name"`
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID      string       `json:"user_id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// AuthUsecase represents the auth usecase contract
type AuthUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
	ValidateToken(token string) (*TokenClaims, error)
	Logout(ctx context.Context, userID string) error
}

// Auth errors
var (
	ErrInvalidCredentials = NewAppError("invalid username or password", 401)
	ErrInvalidToken       = NewAppError("invalid or expired token", 401)
	ErrUserInactive       = NewAppError("user account is inactive", 403)
	ErrUnauthorized       = NewAppError("unauthorized access", 401)
	ErrForbidden          = NewAppError("forbidden: insufficient permissions", 403)
)

// AppError represents application error with status code
type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(message string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
	}
}
