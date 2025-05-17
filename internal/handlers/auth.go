package handlers

import (
	"net/http"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/models"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理程序
type AuthHandler struct {
	userService service.UserService
	jwtService  service.JWTService
}

// NewAuthHandler 创建认证处理程序
func NewAuthHandler(services service.Factory) *AuthHandler {
	return &AuthHandler{
		userService: services.User(),
		jwtService:  services.JWT(),
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c, req.Username, req.Email, req.Password, req.PairToken)
	if err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "用户已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	tokenResp, err := h.generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusCreated, tokenResp)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User
	var err error

	// 根据提供的登录方式选择登录方法
	switch {
	case req.Email != "":
		user, err = h.userService.Login(c, req.Email, req.Password)
	case req.Username != "":
		user, err = h.userService.LoginByUsername(c, req.Username, req.Password)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名或邮箱"})
		return
	}

	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidPassword {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	tokenResp, err := h.generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, tokenResp)
}

// 生成JWT令牌
func (h *AuthHandler) generateTokens(userID int64) (*dto.TokenResponse, error) {
	// 使用JWT服务生成令牌
	tokenDetails, err := h.jwtService.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
		ExpiresIn:    tokenDetails.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从请求中获取刷新令牌
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用JWT服务验证刷新令牌
	_, claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的刷新令牌"})
		return
	}

	// 检查令牌类型
	tokenType, ok := claims["typ"].(string)
	if !ok || tokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌类型"})
		return
	}

	// 获取用户ID
	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户ID"})
		return
	}
	userID := int64(userIDFloat)

	// 生成新令牌
	tokenResp, err := h.generateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, tokenResp)
}
