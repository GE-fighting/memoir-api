package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListMediaHandler lists media items
func ListPhotoVideoHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement list media logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("获取媒体列表成功"))
	}
}

// UploadMediaHandler uploads a new media item
func UploadPhotoVideoHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreatePhotoVideoRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "Invalid request", err.Error()))
			return
		}

		photoVideo, err := services.PhotoVideo().CreatePhotoVideo(c, &request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "Failed to create photo video", err.Error()))
			return
		}

		c.JSON(http.StatusCreated, dto.NewSuccessResponse(photoVideo, "上传媒体成功"))
	}
}
