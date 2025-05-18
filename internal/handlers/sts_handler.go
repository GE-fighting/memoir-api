package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"memoir-api/internal/aliyun"
)

// GenerateSTSToken generates a temporary STS token for OSS access
// @Summary Generate an STS token for Aliyun OSS access
// @Description Generates a temporary STS token for accessing Aliyun OSS
// @Tags oss
// @Accept json
// @Produce json
// @Param category body STSRequest true "Category for uploads (optional)"
// @Success 200 {object} aliyun.STSToken
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oss/token [post]
func GenerateSTSToken(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	// For development, we'll use a hardcoded user ID

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Generate STS token
	token, err := aliyun.GenerateSTSToken(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate STS token: " + err.Error(),
		})
		return
	}

	// Return token
	c.JSON(http.StatusOK, token)
}
