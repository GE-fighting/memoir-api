package api

import (
	"errors"
	"memoir-api/internal/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error codes
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

// Error variables
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

// HandleError function
func HandleError(c *gin.Context, err error) {
	var response dto.Response
	switch {
	case errors.Is(err, ErrInvalidInput):
		response = dto.NewErrorResponse(http.StatusBadRequest, ErrCodeBadRequest, err.Error())
	case errors.Is(err, ErrUnauthorized):
		response = dto.NewErrorResponse(http.StatusUnauthorized, ErrCodeUnauthorized, err.Error())
	case errors.Is(err, ErrForbidden):
		response = dto.NewErrorResponse(http.StatusForbidden, ErrCodeForbidden, err.Error())
	case errors.Is(err, ErrNotFound):
		response = dto.NewErrorResponse(http.StatusNotFound, ErrCodeNotFound, err.Error())
	case errors.Is(err, ErrConflict):
		response = dto.NewErrorResponse(http.StatusConflict, ErrCodeConflict, err.Error())
	case errors.Is(err, ErrValidation):
		response = dto.NewErrorResponse(http.StatusBadRequest, ErrCodeValidation, err.Error())
	case errors.Is(err, ErrResourceUnavailable):
		response = dto.NewErrorResponse(http.StatusServiceUnavailable, ErrCodeResourceUnavailable, err.Error())
	default:
		response = dto.NewErrorResponse(http.StatusInternalServerError, ErrCodeInternalServer, "An unexpected error occurred")
	}
	c.JSON(response.Code, response)
}
