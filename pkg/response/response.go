package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response
type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Meta       *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// Success sends a success response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

// SuccessWithMeta sends a success response with pagination metadata
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta *Meta) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       meta,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       err,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err string) {
	Error(c, http.StatusBadRequest, message, err)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, "resource not found")
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, err string) {
	Error(c, http.StatusInternalServerError, message, err)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string, err string) {
	Error(c, http.StatusConflict, message, err)
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}
