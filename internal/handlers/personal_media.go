package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"memoir-api/internal/service"
)

// PersonalMediaHandler 个人媒体处理器
type PersonalMediaHandler struct {
	personalMediaService service.PersonalMediaService
}

// NewPersonalMediaHandler 创建个人媒体处理器
func NewPersonalMediaHandler(personalMediaService service.PersonalMediaService) *PersonalMediaHandler {
	return &PersonalMediaHandler{
		personalMediaService: personalMediaService,
	}
}

// RegisterRoutes 注册路由
func (h *PersonalMediaHandler) RegisterRoutes(router *gin.Engine) {
	personalGroup := router.Group("/api/personal-media")
	personalGroup.Use(AuthMiddleware()) // 确保用户已认证

	personalGroup.POST("", h.CreatePersonalMedia)
	personalGroup.GET("", h.QueryPersonalMedia)
	personalGroup.GET("/:id", h.GetPersonalMediaByID)
	personalGroup.PUT("/:id", h.UpdatePersonalMedia)
	personalGroup.DELETE("/:id", h.DeletePersonalMedia)
}

// CreatePersonalMediaRequest 创建个人媒体请求
type CreatePersonalMediaRequest struct {
	MediaType string   `form:"mediaType" binding:"required,oneof=photo video"`
	Category  string   `form:"category" binding:"required"`
	Title     string   `form:"title" binding:"required"`
	IsPrivate bool     `form:"isPrivate"`
	Tags      []string `form:"tags"`
}

// CreatePersonalMedia 创建个人媒体
// @Summary 创建个人媒体
// @Description 上传个人照片或视频
// @Tags 个人媒体
// @Accept multipart/form-data
// @Produce json
// @Param mediaType formData string true "媒体类型 (photo/video)"
// @Param category formData string true "分类"
// @Param title formData string true "标题"
// @Param description formData string false "描述(JSON)"
// @Param file formData file true "媒体文件"
// @Param isPrivate formData bool false "是否私密"
// @Param tags formData []string false "标签"
// @Success 201 {object} models.PersonalMedia
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/personal-media [post]
func (h *PersonalMediaHandler) CreatePersonalMedia(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    ErrCodeUnauthorized,
			Message: "未授权",
		})
		return
	}

	var req CreatePersonalMediaRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 获取描述字段（JSON）
	description := []byte(c.PostForm("description"))

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的文件: " + err.Error(),
		})
		return
	}

	// 创建个人媒体
	media, err := h.personalMediaService.Create(c, userID, req.MediaType, req.Category, req.Title, description, file, req.IsPrivate, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    ErrCodeInternalServer,
			Message: "创建个人媒体失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, media)
}

// QueryPersonalMediaRequest 查询个人媒体请求
type QueryPersonalMediaRequest struct {
	Category  string `form:"category"`
	MediaType string `form:"mediaType"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"pageSize,default=20"`
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Data      interface{} `json:"data"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
	TotalPage int64       `json:"totalPage"`
}

// QueryPersonalMedia 查询个人媒体
// @Summary 查询个人媒体
// @Description 根据条件查询个人媒体
// @Tags 个人媒体
// @Accept json
// @Produce json
// @Param category query string false "分类"
// @Param mediaType query string false "媒体类型 (photo/video)"
// @Param page query int false "页码, 默认1"
// @Param pageSize query int false "每页数量, 默认20"
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/personal-media [get]
func (h *PersonalMediaHandler) QueryPersonalMedia(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    ErrCodeUnauthorized,
			Message: "未授权",
		})
		return
	}

	var req QueryPersonalMediaRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 查询个人媒体
	media, total, err := h.personalMediaService.Query(c, userID, req.Category, req.MediaType, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    ErrCodeInternalServer,
			Message: "查询个人媒体失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      media,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	})
}

// GetPersonalMediaByID 获取单个个人媒体
// @Summary 获取单个个人媒体
// @Description 根据ID获取个人媒体
// @Tags 个人媒体
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Success 200 {object} models.PersonalMedia
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "媒体不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/personal-media/{id} [get]
func (h *PersonalMediaHandler) GetPersonalMediaByID(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    ErrCodeUnauthorized,
			Message: "未授权",
		})
		return
	}

	// 获取媒体ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的ID",
		})
		return
	}

	// 获取媒体
	media, err := h.personalMediaService.GetByID(c, id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		errCode := ErrCodeInternalServer
		message := err.Error()

		if errors.Is(err, errors.New("无权访问此媒体")) {
			status = http.StatusForbidden
			errCode = ErrCodeForbidden
		} else if errors.Is(err, errors.New("媒体不存在")) {
			status = http.StatusNotFound
			errCode = ErrCodeNotFound
		}

		c.JSON(status, ErrorResponse{
			Code:    errCode,
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, media)
}

// UpdatePersonalMediaRequest 更新个人媒体请求
type UpdatePersonalMediaRequest struct {
	Title       string          `json:"title" binding:"required"`
	Description json.RawMessage `json:"description"`
	IsPrivate   bool            `json:"isPrivate"`
	Tags        []string        `json:"tags"`
}

// UpdatePersonalMedia 更新个人媒体
// @Summary 更新个人媒体
// @Description 更新个人媒体信息
// @Tags 个人媒体
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Param request body UpdatePersonalMediaRequest true "更新请求"
// @Success 200 {object} models.PersonalMedia
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 403 {object} ErrorResponse "禁止访问"
// @Failure 404 {object} ErrorResponse "媒体不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/personal-media/{id} [put]
func (h *PersonalMediaHandler) UpdatePersonalMedia(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    ErrCodeUnauthorized,
			Message: "未授权",
		})
		return
	}

	// 获取媒体ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的ID",
		})
		return
	}

	var req UpdatePersonalMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 更新媒体
	media, err := h.personalMediaService.Update(c, id, userID, req.Title, req.Description, req.IsPrivate, req.Tags)
	if err != nil {
		status := http.StatusInternalServerError
		errCode := ErrCodeInternalServer
		message := err.Error()

		if errors.Is(err, errors.New("无权访问此媒体")) {
			status = http.StatusForbidden
			errCode = ErrCodeForbidden
		} else if errors.Is(err, errors.New("媒体不存在")) {
			status = http.StatusNotFound
			errCode = ErrCodeNotFound
		}

		c.JSON(status, ErrorResponse{
			Code:    errCode,
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, media)
}

// DeletePersonalMedia 删除个人媒体
// @Summary 删除个人媒体
// @Description 删除个人媒体
// @Tags 个人媒体
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Success 204 {object} nil "删除成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 403 {object} ErrorResponse "禁止访问"
// @Failure 404 {object} ErrorResponse "媒体不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/personal-media/{id} [delete]
func (h *PersonalMediaHandler) DeletePersonalMedia(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    ErrCodeUnauthorized,
			Message: "未授权",
		})
		return
	}

	// 获取媒体ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "无效的ID",
		})
		return
	}

	// 删除媒体
	err = h.personalMediaService.Delete(c, id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		errCode := ErrCodeInternalServer
		message := err.Error()

		if errors.Is(err, errors.New("无权访问此媒体")) {
			status = http.StatusForbidden
			errCode = ErrCodeForbidden
		} else if errors.Is(err, errors.New("媒体不存在")) {
			status = http.StatusNotFound
			errCode = ErrCodeNotFound
		}

		c.JSON(status, ErrorResponse{
			Code:    errCode,
			Message: message,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
