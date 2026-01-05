package http

import (
	"net/http"
	"strconv"

	"icafe-registration/internal/domain"
	"icafe-registration/pkg/response"
	"icafe-registration/pkg/validator"

	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for users
type UserHandler struct {
	userUsecase domain.UserUsecase
	validator   *validator.CustomValidator
}

// NewUserHandler creates a new user handler
func NewUserHandler(router *gin.RouterGroup, uc domain.UserUsecase) {
	handler := &UserHandler{
		userUsecase: uc,
		validator:   validator.NewValidator(),
	}

	router.POST("/users", handler.Create)
	router.GET("/users", handler.GetAll)
	router.GET("/users/:id", handler.GetByID)
	router.PUT("/users/:id", handler.Update)
	router.PUT("/users/:id/password", handler.ChangePassword)
	router.DELETE("/users/:id", handler.Delete)
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body domain.CreateUserRequest true "User data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req domain.CreateUserRequest
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

	user, err := h.userUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrAlreadyExists:
			response.Conflict(c, "Username already exists", err.Error())
		case domain.ErrEmailAlreadyExists:
			response.Conflict(c, "Email already exists", err.Error())
		default:
			response.InternalServerError(c, "Failed to create user", err.Error())
		}
		return
	}

	response.Created(c, "User created successfully", user)
}

// GetAll godoc
// @Summary Get all users
// @Description Get all users with pagination
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	users, total, err := h.userUsecase.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get users", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Users retrieved successfully", users, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetByID godoc
// @Summary Get a user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "User not found")
		default:
			response.InternalServerError(c, "Failed to get user", err.Error())
		}
		return
	}

	response.OK(c, "User retrieved successfully", user)
}

// Update godoc
// @Summary Update a user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body domain.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateUserRequest
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

	user, err := h.userUsecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "User not found")
		case domain.ErrEmailAlreadyExists:
			response.Conflict(c, "Email already exists", err.Error())
		default:
			response.InternalServerError(c, "Failed to update user", err.Error())
		}
		return
	}

	response.OK(c, "User updated successfully", user)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param password body domain.ChangePasswordRequest true "Password data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var req domain.ChangePasswordRequest
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

	err := h.userUsecase.ChangePassword(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "User not found")
		case domain.ErrInvalidCredentials:
			response.BadRequest(c, "Invalid old password", err.Error())
		default:
			response.InternalServerError(c, "Failed to change password", err.Error())
		}
		return
	}

	response.OK(c, "Password changed successfully", nil)
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "User not found")
		default:
			response.InternalServerError(c, "Failed to delete user", err.Error())
		}
		return
	}

	response.OK(c, "User deleted successfully", nil)
}
