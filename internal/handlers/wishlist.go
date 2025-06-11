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
