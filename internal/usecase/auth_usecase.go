package usecase

import (
	"context"
	"time"

	"icafe-registration/internal/config"
	"icafe-registration/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo       domain.UserRepository
	jwtConfig      *config.JWTConfig
	contextTimeout time.Duration
}

// JWTClaims represents the claims in JWT token
type JWTClaims struct {
	UserID      string              `json:"user_id"`
	Username    string              `json:"username"`
	Email       string              `json:"email"`
	Role        domain.Role         `json:"role"`
	Permissions []domain.Permission `json:"permissions"`
	jwt.RegisteredClaims
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo domain.UserRepository, jwtConfig *config.JWTConfig, timeout time.Duration) domain.AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		jwtConfig:      jwtConfig,
		contextTimeout: timeout,
	}
}

// Register creates a new user account (public registration with sale role)
func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Check if username already exists
	existingUser, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrAlreadyExists
	}

	// Check if phone already exists
	existingUser, err = u.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrPhoneAlreadyExists
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user with default sale role
	user := &domain.User{
		Username:    req.Username,
		Phone:       req.Phone,
		Password:    hashedPassword,
		FullName:    req.FullName,
		Role:        domain.RoleSale,
		Permissions: domain.GetPermissionsForRole(domain.RoleSale),
		IsActive:    true,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates user and returns tokens
func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Find user by username
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Update last login
	u.userRepo.UpdateLastLogin(ctx, user.ID.Hex())

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    u.jwtConfig.AccessTokenDuration * 60, // Convert to seconds
		User: &domain.UserInfo{
			ID:          user.ID.Hex(),
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			Role:        user.Role,
			Permissions: user.Permissions,
		},
	}, nil
}

// RefreshToken generates new tokens from refresh token
func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate refresh token
	claims, err := u.validateToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Get user from database
	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// Generate new tokens
	newAccessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    u.jwtConfig.AccessTokenDuration * 60,
		User: &domain.UserInfo{
			ID:          user.ID.Hex(),
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			Role:        user.Role,
			Permissions: user.Permissions,
		},
	}, nil
}

// ValidateToken validates JWT token and returns claims
func (u *authUsecase) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	claims, err := u.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &domain.TokenClaims{
		UserID:      claims.UserID,
		Username:    claims.Username,
		Email:       claims.Email,
		Role:        claims.Role,
		Permissions: claims.Permissions,
	}, nil
}

// Logout handles user logout (can be extended to blacklist tokens)
func (u *authUsecase) Logout(ctx context.Context, userID string) error {
	// For now, just return nil
	// In production, you might want to blacklist the token
	return nil
}

// generateAccessToken generates a new access token
func (u *authUsecase) generateAccessToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(u.jwtConfig.AccessTokenDuration) * time.Minute)

	claims := &JWTClaims{
		UserID:      user.ID.Hex(),
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtConfig.SecretKey))
}

// generateRefreshToken generates a new refresh token
func (u *authUsecase) generateRefreshToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(u.jwtConfig.RefreshTokenDuration) * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtConfig.SecretKey))
}

// validateToken validates a JWT token
func (u *authUsecase) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(u.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
