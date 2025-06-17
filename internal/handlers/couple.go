package handlers

import (
	"fmt"
	"memoir-api/internal/aliyun"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCoupleHandler gets the current user's couple
func CreateCoupleHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")
		var req dto.CreateCoupleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
			return
		}
		req.UserID = userId
		_, err := services.Couple().CreateCouple(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建情侣关系失败", err.Error()))
			return
		}
		coupleInfo, err := services.Couple().GetCoupleInfo(c, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取情侣信息失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(coupleInfo))
	}
}

func GenerateCoupleSTSToken(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		coupleID, err := services.User().GetCoupleID(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取情侣关系失败", err.Error()))
			return
		}
		// Generate STS token
		token, err := aliyun.GenerateSTSToken(c.Request.Context(), fmt.Sprintf("%v", coupleID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "生成STS令牌失败", err.Error()))
			return
		}

		// Return token
		c.JSON(http.StatusOK, dto.NewSuccessResponse(token))
	}
}

func GetCoupleInfoHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")
		coupleInfo, err := services.Couple().GetCoupleInfo(c, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取情侣信息失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(coupleInfo))
	}
}
