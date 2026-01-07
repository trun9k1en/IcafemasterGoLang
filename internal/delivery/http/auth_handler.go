package http

import (
	"icafe-registration/internal/domain"
	"icafe-registration/pkg/response"
	"icafe-registration/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AuthHandler represents the HTTP handler for authentication
type AuthHandler struct {
	authUsecase domain.AuthUsecase
	validator   *validator.CustomValidator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(router *gin.RouterGroup, uc domain.AuthUsecase) {
	handler := &AuthHandler{
		authUsecase: uc,
		validator:   validator.NewValidator(),
	}

	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/logout", handler.Logout)
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with username, password and phone (default role: sale)
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.RegisterRequest true "Registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		errors := validator.GetValidationErrors(err)
		response.BadRequest(c, "Validation failed", mapToString(errors))
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrAlreadyExists:
			response.Conflict(c, "Username already exists", err.Error())
		case domain.ErrPhoneAlreadyExists:
			response.Conflict(c, "Phone number already registered", err.Error())
		default:
			response.InternalServerError(c, "Registration failed", err.Error())
		}
		return
	}

	response.Created(c, "Registration successful", user)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		errors := validator.GetValidationErrors(err)
		response.BadRequest(c, "Validation failed", mapToString(errors))
		return
	}

	loginResponse, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Message, appErr.Message)
			return
		}
		response.InternalServerError(c, "Login failed", err.Error())
		return
	}

	response.OK(c, "Login successful", loginResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		errors := validator.GetValidationErrors(err)
		response.BadRequest(c, "Validation failed", mapToString(errors))
		return
	}

	loginResponse, err := h.authUsecase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Message, appErr.Message)
			return
		}
		response.InternalServerError(c, "Token refresh failed", err.Error())
		return
	}

	response.OK(c, "Token refreshed successfully", loginResponse)
}

// Logout godoc
// @Summary User logout
// @Description Logout user (invalidate token)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.OK(c, "Logged out successfully", nil)
		return
	}

	if err := h.authUsecase.Logout(c.Request.Context(), userID.(string)); err != nil {
		response.InternalServerError(c, "Logout failed", err.Error())
		return
	}

	response.OK(c, "Logged out successfully", nil)
}
