package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListLocationsHandler lists locations
func ListLocationsHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement list locations logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("获取位置列表成功"))
	}
}

// CreateLocationHandler creates a new location
func CreateLocationHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement create location logic
		c.JSON(http.StatusCreated, dto.EmptySuccessResponse("创建位置成功"))
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
		// TODO: Implement delete location logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("删除位置成功"))
	}
}
