package handlers

import (
	"github.com/gin-gonic/gin"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"
)

func GetDashboardDataHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")
		dashboard, err := services.Dashboard().GetDashboardData(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取仪表盘数据失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(dashboard))
	}
}
