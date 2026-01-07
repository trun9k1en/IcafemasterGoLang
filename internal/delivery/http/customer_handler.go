package http

import (
	"net/http"
	"strconv"

	"icafe-registration/internal/domain"
	"icafe-registration/pkg/response"
	"icafe-registration/pkg/validator"

	"github.com/gin-gonic/gin"
)

// CustomerHandler represents the HTTP handler for customer
type CustomerHandler struct {
	customerUsecase domain.CustomerUsecase
	validator       *validator.CustomValidator
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(router *gin.RouterGroup, uc domain.CustomerUsecase) {
	handler := &CustomerHandler{
		customerUsecase: uc,
		validator:       validator.NewValidator(),
	}

	customers := router.Group("/customers")
	{
		// Read operations - accessible by admin and sale
		customers.GET("", handler.GetAll)
		customers.GET("/:id", handler.GetByID)

		// Write operations - accessible by admin only
		adminOnly := customers.Group("")
		adminOnly.Use(RequireRole(domain.RoleAdmin))
		{
			adminOnly.POST("", handler.Create)
			adminOnly.PUT("/:id", handler.Update)
			adminOnly.DELETE("/:id", handler.Delete)
		}
	}
}

// Create godoc
// @Summary Create a new customer
// @Description Create a new customer with the provided data (admin only)
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customer body domain.CreateCustomerRequest true "Customer data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req domain.CreateCustomerRequest
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

	customer, err := h.customerUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrPhoneAlreadyExists:
			response.Conflict(c, "Phone number already registered", err.Error())
		default:
			response.InternalServerError(c, "Failed to create customer", err.Error())
		}
		return
	}

	response.Created(c, "Customer created successfully", customer)
}

// GetAll godoc
// @Summary Get all customers
// @Description Get all customers with pagination
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers [get]
func (h *CustomerHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	customers, total, err := h.customerUsecase.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get customers", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Customers retrieved successfully", customers, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetByID godoc
// @Summary Get a customer by ID
// @Description Get a customer by its ID
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	customer, err := h.customerUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "Customer not found")
		default:
			response.InternalServerError(c, "Failed to get customer", err.Error())
		}
		return
	}

	response.OK(c, "Customer retrieved successfully", customer)
}

// Update godoc
// @Summary Update a customer
// @Description Update a customer by its ID (admin only)
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Param customer body domain.UpdateCustomerRequest true "Customer data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateCustomerRequest
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

	customer, err := h.customerUsecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "Customer not found")
		case domain.ErrPhoneAlreadyExists:
			response.Conflict(c, "Phone number already registered", err.Error())
		default:
			response.InternalServerError(c, "Failed to update customer", err.Error())
		}
		return
	}

	response.OK(c, "Customer updated successfully", customer)
}

// Delete godoc
// @Summary Delete a customer
// @Description Delete a customer by its ID (admin only)
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.customerUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "Customer not found")
		default:
			response.InternalServerError(c, "Failed to delete customer", err.Error())
		}
		return
	}

	response.OK(c, "Customer deleted successfully", nil)
}
