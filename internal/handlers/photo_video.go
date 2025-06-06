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
		var parms *dto.PhotoVideoQueryParams
		if err := c.ShouldBindQuery(&parms); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		pageResult, err := services.PhotoVideo().Query(c.Request.Context(), parms)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "查询失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(pageResult))
	}
}

// UploadMediaHandler uploads a new media item
func CreatePhotoVideoHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreatePhotoVideoRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "Invalid request", err.Error()))
			return
		}
		userID := c.GetInt64("user_id")
		request.UserID = userID
		photoVideo, err := services.PhotoVideo().CreatePhotoVideo(c, &request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "Failed to create photo video", err.Error()))
			return
		}
		c.JSON(http.StatusCreated, dto.NewSuccessResponse(photoVideo))
	}
}
