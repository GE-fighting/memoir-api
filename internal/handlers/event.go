package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListTimelineEventsHandler lists timeline events
func PageTimelineEventsHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TimelineEventQueryParams
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		models, total, err := services.TimelineEvent().ListTimelineEventsByCoupleID(
			c.Request.Context(), req.CoupleID, req.Offset(), req.Limit())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK,
			dto.NewSuccessResponse(dto.NewPageResult(models, total, req.Page, req.PageSize)))

	}
}

// CreateTimelineEventHandler creates a new timeline event
func CreateTimelineEventHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTimelineEventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}
		result, err := services.TimelineEvent().CreateTimelineEvent(c.Request.Context(), &req)
		if err != nil || result == false {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建时间线事件失败", err.Error()))
		}
		c.JSON(http.StatusCreated, dto.EmptySuccessResponse("创建时间线事件成功"))
	}
}

// GetTimelineEventHandler gets a specific timeline event
func GetTimelineEventHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		coupleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "获取故事ID失败", err.Error()))
		}
		event, err := services.TimelineEvent().GetTimelineEventByID(c.Request.Context(), coupleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "<UNK>", err.Error()))
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(event))
	}
}

// UpdateTimelineEventHandler updates a timeline event
func UpdateTimelineEventHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update timeline event logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("更新时间线事件成功"))
	}
}

// DeleteTimelineEventHandler deletes a timeline event
func DeleteTimelineEventHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement delete timeline event logic
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("删除时间线事件成功"))
	}
}
