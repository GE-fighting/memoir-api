package handlers

import (
	"fmt"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCoupleAlbumHandler 创建情侣相册
func CreateCoupleAlbumHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateCoupleAlbumRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}

		// 如果请求中没有指定用户ID，则使用当前登录用户的ID
		if req.UserID == 0 {
			userIDValue, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(http.StatusUnauthorized, "用户未登录", "未找到用户ID"))
				return
			}

			userID, ok := userIDValue.(int64)
			if !ok {
				c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "用户ID类型无效", "用户ID类型断言失败"))
				return
			}
			req.UserID = userID
		}

		// 创建相册
		album, err := services.CoupleAlbum().Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建相册失败", err.Error()))
			return
		}

		c.JSON(http.StatusCreated, dto.NewSuccessResponse(album))
	}
}

// GetCoupleAlbumHandler 获取单个情侣相册
func GetCoupleAlbumHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析相册ID
		albumID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的相册ID", err.Error()))
			return
		}

		// 获取相册
		album, err := services.CoupleAlbum().GetByID(c.Request.Context(), albumID)
		if err != nil {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "相册不存在", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(album))
	}
}

// GetCoupleAlbumWithPhotosHandler 获取情侣相册及其照片
func GetCoupleAlbumWithPhotosHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析相册ID
		albumID, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的相册ID", err.Error()))
			return
		}

		// 获取相册及其照片
		album, err := services.CoupleAlbum().GetWithPhotos(c.Request.Context(), albumID)
		if err != nil {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "相册不存在", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(album))
	}
}

// ListCoupleAlbumsHandler 获取情侣的所有相册
func ListCoupleAlbumsHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID（JWT中间件设置）
		userID := c.GetInt64("user_id")
		// 获取用户信息
		user, err := services.User().GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取用户信息失败", err.Error()))
			return
		}

		// 确保用户属于一个情侣关系
		if user.CoupleID == 0 {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "用户不属于任何情侣关系", ""))
			return
		}

		// 获取情侣相册列表
		albums, err := services.CoupleAlbum().GetByCoupleID(c.Request.Context(), user.CoupleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取相册列表失败", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(albums))
	}
}

// UpdateCoupleAlbumHandler 更新情侣相册
func UpdateCoupleAlbumHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析相册ID
		albumID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的相册ID", err.Error()))
			return
		}

		var req dto.UpdateCoupleAlbumRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}

		// 更新相册
		album, err := services.CoupleAlbum().Update(c.Request.Context(), albumID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "更新相册失败", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(album))
	}
}

// DeleteCoupleAlbumHandler 删除情侣相册
func DeleteCoupleAlbumHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析相册ID
		albumID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的相册ID", err.Error()))
			return
		}

		// 删除相册
		if err := services.CoupleAlbum().Delete(c.Request.Context(), albumID); err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "删除相册失败", err.Error()))
			return
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
	}
}

func DeleteCoupleAlbumPhotosHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var delReq dto.DeleteCoupleAlbumPhotosRequest
		if err := c.ShouldBindJSON(&delReq); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}
		// 删除相册照片
		if err := services.PhotoVideo().BatchDeletePhotoVideo(c.Request.Context(), delReq.PhotoVideoIDs); err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "删除相册照片失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
	}
}
