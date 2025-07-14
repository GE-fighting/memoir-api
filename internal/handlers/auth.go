package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/logger"
	"memoir-api/internal/models"
	"memoir-api/internal/service"
)

// AuthHandler 认证处理程序
type AuthHandler struct {
	userService service.UserService
	jwtService  service.JWTService
	emailSvc    service.EmailService
}

// NewAuthHandler 创建认证处理程序
func NewAuthHandler(services service.Factory) *AuthHandler {
	return &AuthHandler{
		userService: services.User(),
		jwtService:  services.JWT(),
		emailSvc:    services.Email(),
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
		return
	}

	user, err := h.userService.Register(c, req.Username, req.Email, req.Password, req.PairToken)
	if err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, dto.NewErrorResponse(http.StatusConflict, "用户已存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "注册失败", err.Error()))
		return
	}

	// 生成JWT令牌
	tokenResp, err := h.generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "生成令牌失败", err.Error()))
		return
	}

	// 验证成功后发送欢迎邮件
	err = h.emailSvc.SendWelcomeEmail(c, req.Email, user.Username)
	if err != nil {
		// 记录错误但不影响验证流程
		logger.GetLogger("auth").Error(err, "发送欢迎邮件错误")
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(tokenResp))
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
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
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请提供用户名或邮箱", "missing login credentials"))
		return
	}

	if err != nil {
		status := http.StatusInternalServerError
		message := "登录失败"
		if err == service.ErrInvalidPassword {
			status = http.StatusUnauthorized
			message = "密码错误"
		}
		c.JSON(status, dto.NewErrorResponse(status, message, err.Error()))
		return
	}

	// 生成JWT令牌
	tokenResp, err := h.generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "生成令牌失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(tokenResp))
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
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
		return
	}

	// 使用JWT服务验证刷新令牌
	_, claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(http.StatusUnauthorized, "无效的刷新令牌", err.Error()))
		return
	}

	// 检查令牌类型
	tokenType, ok := claims["typ"].(string)
	if !ok || tokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(http.StatusUnauthorized, "无效的令牌类型", "invalid token type"))
		return
	}

	// 获取用户ID
	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "无法获取用户ID", "invalid user id in token"))
		return
	}
	userID := int64(userIDFloat)

	// 生成新令牌
	tokenResp, err := h.generateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "生成令牌失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(tokenResp))
}
