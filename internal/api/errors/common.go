package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Status  int    `json:"-"`                 // HTTP status code
	Code    string `json:"code"`              // Error code for clients
	Message string `json:"message"`           // User-friendly error message
	Details any    `json:"details,omitempty"` // Optional additional details
}

// Common error codes
const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeValidation          = "VALIDATION_ERROR"
	ErrCodeResourceUnavailable = "RESOURCE_UNAVAILABLE"
)

// Common application errors
var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource conflict")
	ErrInternalServer      = errors.New("internal server error")
	ErrValidation          = errors.New("validation error")
	ErrResourceUnavailable = errors.New("resource unavailable")
)

// NewErrorResponse creates a new error response with the given parameters
func NewErrorResponse(status int, code, message string, details any) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// HandleError sends an appropriate error response based on the error type
func HandleError(c *gin.Context, err error, details ...any) {
	var detail any
	if len(details) > 0 {
		detail = details[0]
	}

	var response ErrorResponse

	switch {
	case errors.Is(err, ErrInvalidInput):
		response = NewErrorResponse(http.StatusBadRequest, ErrCodeBadRequest, err.Error(), detail)
	case errors.Is(err, ErrUnauthorized):
		response = NewErrorResponse(http.StatusUnauthorized, ErrCodeUnauthorized, err.Error(), detail)
	case errors.Is(err, ErrForbidden):
		response = NewErrorResponse(http.StatusForbidden, ErrCodeForbidden, err.Error(), detail)
	case errors.Is(err, ErrNotFound):
		response = NewErrorResponse(http.StatusNotFound, ErrCodeNotFound, err.Error(), detail)
	case errors.Is(err, ErrConflict):
		response = NewErrorResponse(http.StatusConflict, ErrCodeConflict, err.Error(), detail)
	case errors.Is(err, ErrValidation):
		response = NewErrorResponse(http.StatusBadRequest, ErrCodeValidation, err.Error(), detail)
	case errors.Is(err, ErrResourceUnavailable):
		response = NewErrorResponse(http.StatusServiceUnavailable, ErrCodeResourceUnavailable, err.Error(), detail)
	default:
		// Log unexpected errors but don't expose details to clients
		response = NewErrorResponse(http.StatusInternalServerError, ErrCodeInternalServer, "An unexpected error occurred", nil)
	}

	c.JSON(response.Status, response)
}
