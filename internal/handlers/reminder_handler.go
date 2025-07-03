package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TriggerAnniversaryRemindersHandler 触发纪念日提醒
func TriggerAnniversaryRemindersHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行纪念日提醒
		err := services.CoupleReminder().CheckAndSendAnniversaryReminders(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "执行纪念日提醒失败", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(gin.H{
			"message": "纪念日提醒已触发",
		}))
	}
}

// TriggerFestivalRemindersHandler 触发节日提醒
func TriggerFestivalRemindersHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行节日提醒
		err := services.CoupleReminder().CheckAndSendFestivalReminders(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "执行节日提醒失败", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(gin.H{
			"message": "节日提醒已触发",
		}))
	}
}
