package handlers

import (
	"net/http"
	"strconv"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateAttachmentHandler 创建附件
func CreateAttachmentHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 解析请求
		var req dto.CreateAttachmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}
		// 设置用户ID
		req.UserID = c.GetInt64("user_id")
		// 创建附件
		attachment, err := services.Attachment().CreateAttachment(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建附件失败", err.Error()))
			return
		}

		// 返回创建的附件
		c.JSON(http.StatusCreated, dto.NewSuccessResponse(dto.AttachmentFromModel(attachment)))
	}
}

// GetAttachmentHandler 获取单个附件
func GetAttachmentHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析附件ID
		attachmentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的附件ID", err.Error()))
			return
		}

		// 获取附件
		attachment, err := services.Attachment().GetAttachmentByID(c.Request.Context(), attachmentID)
		if err != nil {
			if err == service.ErrAttachmentNotFound {
				c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "附件不存在", err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取附件失败", err.Error()))
			return
		}

		// 返回附件信息
		c.JSON(http.StatusOK, dto.NewSuccessResponse(dto.AttachmentFromModel(attachment)))
	}
}

// ListAttachmentsHandler 获取附件列表（分页）
func ListAttachmentsHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析查询参数
		var params dto.AttachmentQueryParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}

		// 获取用户ID
		userID := c.GetInt64("user_id")

		// 只能查看自己的附件或所属情侣关系的附件
		user, err := services.User().GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取用户信息失败", err.Error()))
			return
		}

		// 如果没有指定用户ID，则默认查询当前用户的附件
		if params.UserID == 0 {
			params.UserID = userID
		}

		// 如果查询其他用户的附件，必须是同一情侣关系下的用户
		if params.UserID != userID && (user.CoupleID == 0 || params.CoupleID != user.CoupleID) {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(http.StatusForbidden, "无权访问其他用户的附件", ""))
			return
		}

		// 查询附件
		pageResult, err := services.Attachment().QueryAttachments(c.Request.Context(), &params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "查询附件失败", err.Error()))
			return
		}

		// 返回结果
		c.JSON(http.StatusOK, dto.NewSuccessResponse(pageResult))
	}
}

// DeleteAttachmentHandler 删除附件
func DeleteAttachmentHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析附件ID
		attachmentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的附件ID", err.Error()))
			return
		}

		// 获取附件信息，检查权限
		attachment, err := services.Attachment().GetAttachmentByID(c.Request.Context(), attachmentID)
		if err != nil {
			if err == service.ErrAttachmentNotFound {
				c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "附件不存在", err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取附件信息失败", err.Error()))
			return
		}

		// 只有附件所有者或同一情侣关系的用户可以删除
		userID := c.GetInt64("user_id")
		if attachment.UserID != userID {
			// 如果不是自己的附件，检查是否在同一情侣关系
			user, err := services.User().GetUserByID(c.Request.Context(), userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取用户信息失败", err.Error()))
				return
			}

			// 如果不是同一情侣关系，或者附件不属于情侣空间，则无权删除
			if user.CoupleID == 0 || user.CoupleID != attachment.CoupleID || attachment.SpaceType != "couple" {
				c.JSON(http.StatusForbidden, dto.NewErrorResponse(http.StatusForbidden, "无权删除该附件", ""))
				return
			}
		}

		// 删除附件
		if err := services.Attachment().DeleteAttachment(c.Request.Context(), attachmentID); err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "删除附件失败", err.Error()))
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("附件删除成功"))
	}
}
