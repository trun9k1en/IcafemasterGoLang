package http

import (
	"net/http"
	"strconv"

	"icafe-registration/internal/domain"
	"icafe-registration/pkg/response"
	"icafe-registration/pkg/validator"

	"github.com/gin-gonic/gin"
)

// RegistrationHandler represents the HTTP handler for registration
type RegistrationHandler struct {
	registrationUsecase domain.RegistrationUsecase
	validator           *validator.CustomValidator
}

// NewRegistrationHandler creates a new registration handler
func NewRegistrationHandler(router *gin.RouterGroup, uc domain.RegistrationUsecase) {
	handler := &RegistrationHandler{
		registrationUsecase: uc,
		validator:           validator.NewValidator(),
	}

	router.POST("/registrations", handler.Create)
	router.GET("/registrations", handler.GetAll)
	router.GET("/registrations/:id", handler.GetByID)
	router.PUT("/registrations/:id", handler.Update)
	router.DELETE("/registrations/:id", handler.Delete)
}

// Create godoc
// @Summary Create a new registration
// @Description Create a new registration with the provided data
// @Tags registrations
// @Accept json
// @Produce json
// @Param registration body domain.CreateRegistrationRequest true "Registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /registrations [post]
func (h *RegistrationHandler) Create(c *gin.Context) {
	var req domain.CreateRegistrationRequest
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

	registration, err := h.registrationUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists:
			response.Conflict(c, "Email already registered", err.Error())
		default:
			response.InternalServerError(c, "Failed to create registration", err.Error())
		}
		return
	}

	response.Created(c, "Registration created successfully", registration)
}

// GetAll godoc
// @Summary Get all registrations
// @Description Get all registrations with pagination
// @Tags registrations
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /registrations [get]
func (h *RegistrationHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	registrations, total, err := h.registrationUsecase.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get registrations", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Registrations retrieved successfully", registrations, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetByID godoc
// @Summary Get a registration by ID
// @Description Get a registration by its ID
// @Tags registrations
// @Produce json
// @Param id path string true "Registration ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /registrations/{id} [get]
func (h *RegistrationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	registration, err := h.registrationUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "Registration not found")
		default:
			response.InternalServerError(c, "Failed to get registration", err.Error())
		}
		return
	}

	response.OK(c, "Registration retrieved successfully", registration)
}

// Update godoc
// @Summary Update a registration
// @Description Update a registration by its ID
// @Tags registrations
// @Accept json
// @Produce json
// @Param id path string true "Registration ID"
// @Param registration body domain.UpdateRegistrationRequest true "Registration data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /registrations/{id} [put]
func (h *RegistrationHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateRegistrationRequest
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

	registration, err := h.registrationUsecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "Registration not found")
		case domain.ErrEmailAlreadyExists:
			response.Conflict(c, "Email already registered", err.Error())
		default:
			response.InternalServerError(c, "Failed to update registration", err.Error())
		}
		return
	}

	response.OK(c, "Registration updated successfully", registration)
}

// Delete godoc
// @Summary Delete a registration
// @Description Delete a registration by its ID
// @Tags registrations
// @Produce json
// @Param id path string true "Registration ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /registrations/{id} [delete]
func (h *RegistrationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.registrationUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "Registration not found")
		default:
			response.InternalServerError(c, "Failed to delete registration", err.Error())
		}
		return
	}

	response.OK(c, "Registration deleted successfully", nil)
}

// mapToString converts a map to a string for error display
func mapToString(m map[string]string) string {
	result := ""
	for k, v := range m {
		if result != "" {
			result += ", "
		}
		result += k + ": " + v
	}
	return result
}
