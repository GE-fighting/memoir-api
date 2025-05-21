package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListMediaHandler lists media items
func ListMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement list media logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("获取媒体列表成功"))
	}
}

// UploadMediaHandler uploads a new media item
func UploadMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement upload media logic
		c.JSON(http.StatusCreated, dto.EmptySuccessResponse("上传媒体成功"))
	}
}

// GetMediaHandler gets a specific media item
func GetMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get media logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("获取媒体成功"))
	}
}

// UpdateMediaHandler updates a media item
func UpdateMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update media logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("更新媒体成功"))
	}
}

// DeleteMediaHandler deletes a media item
func DeleteMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement delete media logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("删除媒体成功"))
	}
}
