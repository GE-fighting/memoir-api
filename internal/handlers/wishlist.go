package handlers

import (
	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListWishlistItemsHandler lists wishlist items
func ListWishlistItemsHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		coupleId, err := strconv.ParseInt(c.Query("couple_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的情侣ID", err.Error()))
			return
		}
		wishListDTO, err := services.Wishlist().ListWishlistsByCoupleID(c.Request.Context(), coupleId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "查询心愿单失败", err.Error()))
		}
		c.JSON(http.StatusOK, dto.NewSuccessResponse(wishListDTO))
	}
}

// CreateWishlistItemHandler creates a new wishlist item
func CreateWishlistItemHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateWishlistRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}
		wishlist, err := services.Wishlist().CreateWishlist(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "创建心愿单项目失败", err.Error()))
			return
		}
		c.JSON(http.StatusCreated, dto.NewSuccessResponse(wishlist))
	}
}

// UpdateWishlistItemHandler updates a wishlist item
func UpdateWishlistItemHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析请求体
		var req dto.UpdateWishlistRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}

		// 调用服务层方法处理更新逻辑
		updatedWishlist, err := services.Wishlist().UpdateWishlistByRequest(c.Request.Context(), &req)
		if err != nil {
			if err == service.ErrWishlistNotFound {
				c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "心愿不存在", err.Error()))
			} else {
				c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "更新心愿失败", err.Error()))
			}
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, dto.NewSuccessResponse(dto.WishlistFromModel(updatedWishlist)))
	}
}

// UpdateWishlistItemStatusHandler updates the status of a wishlist item
func UpdateWishlistItemStatusHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析路径参数中的ID
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的心愿ID", err.Error()))
			return
		}

		// 解析请求体
		var req dto.UpdateWishlistStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求参数无效", err.Error()))
			return
		}

		// 调用服务更新状态
		err = services.Wishlist().UpdateWishlistStatus(c.Request.Context(), id, req.Status)
		if err != nil {
			if err == service.ErrWishlistNotFound {
				c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "心愿不存在", err.Error()))
			} else {
				c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "更新心愿状态失败", err.Error()))
			}
			return
		}

		// 获取更新后的心愿项目
		wishlist, err := services.Wishlist().GetWishlistByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "获取更新后的心愿失败", err.Error()))
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, dto.NewSuccessResponse(dto.WishlistFromModel(wishlist)))
	}
}

// DeleteWishlistItemHandler deletes a wishlist item
func DeleteWishlistItemHandler(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析路径参数中的ID
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "无效的心愿ID", err.Error()))
			return
		}

		// 调用服务删除心愿项
		err = services.Wishlist().DeleteWishlist(c.Request.Context(), id)
		if err != nil {
			if err == service.ErrWishlistNotFound {
				c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "心愿不存在", err.Error()))
			} else {
				c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "删除心愿失败", err.Error()))
			}
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("心愿已成功删除"))
	}
}
