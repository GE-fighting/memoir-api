package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"memoir-api/internal/aliyun"
	"memoir-api/internal/api/dto"
)

// GenerateSTSToken generates a temporary STS token for OSS access
func GenerateSTSToken(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	// For development, we'll use a hardcoded user ID

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "用户ID必填", "User ID is required"))
		return
	}

	// Generate STS token
	token, err := aliyun.GenerateSTSToken(c.Request.Context(), fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "生成STS令牌失败", err.Error()))
		return
	}

	// Return token
	c.JSON(http.StatusOK, dto.NewSuccessResponse(token, "STS令牌生成成功"))
}
