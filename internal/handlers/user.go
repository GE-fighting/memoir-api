package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCurrentUserHandler gets the current authenticated user
func GetCurrentUserHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID（JWT中间件设置）
		userIDValue, exists := c.Get("userID")
		if !exists {
			response := dto.NewErrorResponse(http.StatusUnauthorized, "用户未登录", "未找到用户ID")
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		userID, ok := userIDValue.(int64)
		if !ok {
			response := dto.NewErrorResponse(http.StatusInternalServerError, "用户ID类型无效", "用户ID类型断言失败")
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// 获取用户服务并查询用户信息
		userService := services.User()
		user, err := userService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			response := dto.NewErrorResponse(http.StatusInternalServerError, "获取用户信息失败", err.Error())
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// 返回用户信息
		response := dto.NewSuccessResponse(user, "获取当前用户成功")
		c.JSON(http.StatusOK, response)
	}
}

// UpdateUserHandler updates user information
func UpdateUserHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update user logic
		response := dto.EmptySuccessResponse("Update user endpoint")
		c.JSON(http.StatusOK, response)
	}
}

// UpdateUserPreferencesHandler updates user preferences
func UpdateUserPreferencesHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update user preferences logic
		response := dto.EmptySuccessResponse("Update user preferences endpoint")
		c.JSON(http.StatusOK, response)
	}
}
