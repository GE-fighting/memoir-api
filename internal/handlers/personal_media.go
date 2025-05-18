package handlers

import (
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

// CreatePersonalMediaWithURLRequest 通过URL创建个人媒体请求
type CreatePersonalMediaWithURLRequest struct {
	MediaType    string   `json:"mediaType" binding:"required,oneof=photo video"`
	Category     string   `json:"category" binding:"required"`
	Title        string   `json:"title" binding:"required"`
	MediaURL     string   `json:"mediaUrl" binding:"required"`
	ThumbnailURL string   `json:"thumbnailUrl"`
	Description  string   `json:"description"` // JSON字符串
	IsPrivate    bool     `json:"isPrivate"`
	Tags         []string `json:"tags"`
}

// CreatePersonalMediaWithURL 通过URL创建个人媒体（前端直接上传到OSS）
// @Summary 通过URL创建个人媒体
// @Description 前端上传到OSS后，保存媒体信息
// @Tags 个人媒体
// @Accept json
// @Produce json
// @Param request body CreatePersonalMediaWithURLRequest true "创建媒体请求"
// @Success 201 {object} models.PersonalMedia
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/v1/personal-media/url [post]
func CreatePersonalMediaWithURLHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    ErrCodeUnauthorized,
				Message: "未授权",
			})
			return
		}

		var req CreatePersonalMediaWithURLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "无效的请求参数: " + err.Error(),
			})
			return
		}

		// 转换描述为JSON字节
		description := []byte(req.Description)

		// 创建个人媒体
		media, err := services.PersonalMedia().CreateWithURL(
			c,
			userID.(int64),
			req.MediaType,
			req.Category,
			req.Title,
			description,
			req.MediaURL,
			req.ThumbnailURL,
			req.IsPrivate,
			req.Tags,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    ErrCodeInternalServer,
				Message: "创建个人媒体失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, media)
	}
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

// QueryPersonalMediaHandler 查询个人媒体处理函数
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
// @Router /api/v1/personal-media [get]
func QueryPersonalMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
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

		// 查询个人媒体（从数据库获取，而非直接查询OSS）
		media, total, err := services.PersonalMedia().Query(c, userID.(int64), req.Category, req.MediaType, req.Page, req.PageSize)
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
}

// GetPersonalMediaByIDHandler 获取单个个人媒体处理函数
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
// @Router /api/v1/personal-media/{id} [get]
func GetPersonalMediaByIDHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
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
		media, err := services.PersonalMedia().GetByID(c, id, userID.(int64))
		if err != nil {
			// 根据错误类型返回不同状态码
			switch err.Error() {
			case "无权访问此媒体":
				c.JSON(http.StatusForbidden, ErrorResponse{
					Code:    "FORBIDDEN",
					Message: err.Error(),
				})
			case "媒体不存在":
				c.JSON(http.StatusNotFound, ErrorResponse{
					Code:    "NOT_FOUND",
					Message: err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, media)
	}
}

// UpdatePersonalMediaRequest 更新个人媒体请求
type UpdatePersonalMediaRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	IsPrivate   bool     `json:"isPrivate"`
	Tags        []string `json:"tags"`
}

// UpdatePersonalMediaHandler 更新个人媒体处理函数
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
// @Router /api/v1/personal-media/{id} [put]
func UpdatePersonalMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
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

		// 转换描述为JSON字节
		description := []byte(req.Description)

		// 更新媒体
		media, err := services.PersonalMedia().Update(c, id, userID.(int64), req.Title, description, req.IsPrivate, req.Tags)
		if err != nil {
			// 根据错误类型返回不同状态码
			switch err.Error() {
			case "无权访问此媒体":
				c.JSON(http.StatusForbidden, ErrorResponse{
					Code:    "FORBIDDEN",
					Message: err.Error(),
				})
			case "媒体不存在":
				c.JSON(http.StatusNotFound, ErrorResponse{
					Code:    "NOT_FOUND",
					Message: err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, media)
	}
}

// DeletePersonalMediaHandler 删除个人媒体处理函数
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
// @Router /api/v1/personal-media/{id} [delete]
func DeletePersonalMediaHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
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
		err = services.PersonalMedia().Delete(c, id, userID.(int64))
		if err != nil {
			// 根据错误类型返回不同状态码
			switch err.Error() {
			case "无权访问此媒体":
				c.JSON(http.StatusForbidden, ErrorResponse{
					Code:    "FORBIDDEN",
					Message: err.Error(),
				})
			case "媒体不存在":
				c.JSON(http.StatusNotFound, ErrorResponse{
					Code:    "NOT_FOUND",
					Message: err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: err.Error(),
				})
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}
