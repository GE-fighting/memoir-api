package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
)

// CreatePersonalMediaWithURL 通过URL创建个人媒体（前端直接上传到OSS）
func CreatePersonalMediaWithURLHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(http.StatusUnauthorized, "未授权", "用户ID不存在"))
			return
		}

		var req dto.CreatePersonalMediaWithURLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的请求参数", err.Error()))
			return
		}
		req.UserID = userID.(int64)
		// 创建个人媒体
		media, err := services.PersonalMedia().CreateWithURL(
			c,
			req,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建个人媒体失败", err.Error()))
			return
		}

		c.JSON(http.StatusCreated, dto.NewSuccessResponse(media, "创建个人媒体成功"))
	}
}

func PageQueryPersonalMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(http.StatusUnauthorized, "未授权", "用户ID不存在"))
			return
		}
		var req dto.QueryPersonalMediaRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的请求参数", err.Error()))
			return
		}

		// 检查并设置默认值
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}

		req.UserID = userID.(int64)
		media, total, err := services.PersonalMedia().PageQuery(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "查询个人媒体失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(
			dto.NewPaginatedResponse(media, total, req.Page, req.PageSize),
			"查询个人媒体成功",
		))
	}
}
