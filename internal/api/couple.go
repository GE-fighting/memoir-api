package api

import (
	"net/http"
	"strconv"

	"memoir-api/internal/models"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
)

// CoupleHandler 情侣关系处理程序
type CoupleHandler struct {
	coupleService service.CoupleService
}

// NewCoupleHandler 创建情侣关系处理程序
func NewCoupleHandler(coupleService service.CoupleService) *CoupleHandler {
	return &CoupleHandler{
		coupleService: coupleService,
	}
}

// CreateCoupleRequest 创建情侣关系请求
type CreateCoupleRequest struct {
	TimelinePrivacy       string `json:"timeline_privacy" binding:"omitempty,oneof=public private"`
	AlbumPrivacy          string `json:"album_privacy" binding:"omitempty,oneof=public private"`
	AutoGenerateVideo     *bool  `json:"auto_generate_video" binding:"omitempty"`
	ReminderNotifications *bool  `json:"reminder_notifications" binding:"omitempty"`
}

// CreateCouple 创建情侣关系
func (h *CoupleHandler) CreateCouple(c *gin.Context) {
	var req CreateCoupleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	couple := &models.Couple{
		TimelinePrivacy: req.TimelinePrivacy,
		AlbumPrivacy:    req.AlbumPrivacy,
	}

	if req.AutoGenerateVideo != nil {
		couple.AutoGenerateVideo = *req.AutoGenerateVideo
	}

	if req.ReminderNotifications != nil {
		couple.ReminderNotifications = *req.ReminderNotifications
	}

	createdCouple, err := h.coupleService.CreateCouple(c, couple)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建情侣关系失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"couple": createdCouple,
	})
}

// GetCouple 获取情侣关系
func (h *CoupleHandler) GetCouple(c *gin.Context) {
	// 从路径参数获取情侣ID
	coupleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的情侣ID"})
		return
	}

	couple, err := h.coupleService.GetCoupleByID(c, coupleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "情侣关系不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"couple": couple,
	})
}

// UpdateCoupleRequest 更新情侣关系请求
type UpdateCoupleRequest struct {
	TimelinePrivacy       string `json:"timeline_privacy" binding:"omitempty,oneof=public private"`
	AlbumPrivacy          string `json:"album_privacy" binding:"omitempty,oneof=public private"`
	AutoGenerateVideo     *bool  `json:"auto_generate_video" binding:"omitempty"`
	ReminderNotifications *bool  `json:"reminder_notifications" binding:"omitempty"`
}

// UpdateCouple 更新情侣关系
func (h *CoupleHandler) UpdateCouple(c *gin.Context) {
	// 从路径参数获取情侣ID
	coupleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的情侣ID"})
		return
	}

	// 获取当前情侣关系
	couple, err := h.coupleService.GetCoupleByID(c, coupleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "情侣关系不存在"})
		return
	}

	// 解析请求
	var req UpdateCoupleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新情侣关系
	if req.TimelinePrivacy != "" {
		couple.TimelinePrivacy = req.TimelinePrivacy
	}
	if req.AlbumPrivacy != "" {
		couple.AlbumPrivacy = req.AlbumPrivacy
	}
	if req.AutoGenerateVideo != nil {
		couple.AutoGenerateVideo = *req.AutoGenerateVideo
	}
	if req.ReminderNotifications != nil {
		couple.ReminderNotifications = *req.ReminderNotifications
	}

	// 保存更新后的情侣关系
	if err := h.coupleService.UpdateCouple(c, couple); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新情侣关系失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"couple": couple,
	})
}

// DeleteCouple 删除情侣关系
func (h *CoupleHandler) DeleteCouple(c *gin.Context) {
	// 从路径参数获取情侣ID
	coupleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的情侣ID"})
		return
	}

	// 删除情侣关系（这将同时删除所有关联的用户）
	if err := h.coupleService.DeleteCouple(c, coupleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除情侣关系失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "情侣关系已删除"})
}

// GetCoupleUsers 获取情侣关系下的用户
func (h *CoupleHandler) GetCoupleUsers(c *gin.Context) {
	// 从路径参数获取情侣ID
	coupleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的情侣ID"})
		return
	}

	users, err := h.coupleService.GetCoupleUsers(c, coupleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
