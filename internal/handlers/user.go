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
		userIDValue, exists := c.Get("user_id")
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
		response := dto.NewSuccessResponse(user)
		c.JSON(http.StatusOK, response)
	}
}

func ExistCoupleHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		exist, err := services.User().ExistCouple(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取情侣关系失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(exist))
	}
}

func UpdateUserHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		var updateUserRequest dto.UpdateUserRequest
		if c.ShouldBind(&updateUserRequest) != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "参数绑定失败", "请检查参数"))
			return
		}
		updateUserRequest.UserID = userID
		userService := services.User()
		err := userService.UpdateUser(c.Request.Context(), &updateUserRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "更新用户信息失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(updateUserRequest))
	}
}

func UpdatePassword(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")
		var updateDTO dto.UpdateUserPasswordDTO
		if c.ShouldBindJSON(&updateDTO) != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "参数绑定失败", "请检查参数"))
			return

		}
		updateDTO.UserID = userId
		err := services.User().UpdatePassword(c.Request.Context(), &updateDTO)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "更新密码失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("更新成功"))
	}
}
