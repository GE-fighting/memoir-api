package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/logger"
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
		var req dto.UpdateTimelineEventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}

		// 先获取现有事件
		existingEvent, err := services.TimelineEvent().GetTimelineEventByID(c.Request.Context(), req.EventId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取时间线事件失败", err.Error()))
			return
		}

		// 应用更新到已有事件上
		err = req.ApplyToModel(existingEvent)
		if err != nil {
			log := logger.FromContext(c).WithComponent("event_handler")
			log.Error(err, "应用更新错误")
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "应用更新数据失败", err.Error()))
			return
		}

		_, err = services.TimelineEvent().UpdateTimelineEvent(c.Request.Context(), existingEvent, req.LocationIDs, req.PhotoVideoIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "更新回忆失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("更新回忆成功"))
	}
}

// DeleteTimelineEventHandler deletes a timeline event
func DeleteTimelineEventHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		evevtId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "获取时间线事件ID失败", err.Error()))
		}
		err = services.TimelineEvent().DeleteTimelineEvent(c.Request.Context(), evevtId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "删除时间线事件失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("删除时间线事件成功"))
	}
}
