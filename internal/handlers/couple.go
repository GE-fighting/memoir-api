package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCoupleHandler gets the current user's couple
func CreateCoupleHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(http.StatusUnauthorized, "未授权", "用户ID不存在"))
			return
		}
		var req dto.CreateCoupleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
			return
		}
		couple, err := services.Couple().CreateCouple(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建情侣关系失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(couple, "创建情侣关系成功"))
	}
}
