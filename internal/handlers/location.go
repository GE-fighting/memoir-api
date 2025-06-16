package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListLocationsHandler lists locations
func ListLocationsHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, err := strconv.ParseInt(c.Query("couple_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的情侣ID", err.Error()))
			return
		}
		var req dto.LocationQueryParams
		locations, _, err := services.Location().ListLocationsByCoupleID(c.Request.Context(), value, 0, 0)
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "查询所有的地点出错", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(locations))
	}
}

// CreateLocationHandler creates a new location
func CreateLocationHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateLocationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}
		location, err := services.Location().CreateLocation(c.Request.Context(), req.ToModel())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建位置失败", err.Error()))
			return
		}
		c.JSON(http.StatusCreated, dto.NewSuccessResponse(location))
	}
}

// GetLocationHandler gets a specific location
func GetLocationHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get location logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("获取位置成功"))
	}
}

// UpdateLocationHandler updates a location
func UpdateLocationHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update location logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("更新位置成功"))
	}
}

// DeleteLocationHandler deletes a location
func DeleteLocationHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		locationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "地点ID出错", err.Error()))
		}
		err = services.Location().DeleteLocation(c.Request.Context(), locationId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "删除地点出错", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("删除地点成功"))

	}
}
