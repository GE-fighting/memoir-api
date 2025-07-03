package handlers

import (
	"net/http"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
)

// EmailHandler 邮件相关处理器
type EmailHandler struct {
	userService  service.UserService
	emailService service.EmailService
}

// NewEmailHandler 创建邮件处理器
func NewEmailHandler(services service.Factory) *EmailHandler {
	return &EmailHandler{
		userService:  services.User(),
		emailService: services.Email(),
	}
}

// VerifyEmail 处理邮箱验证请求
func (h *EmailHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
		return
	}

	err := h.userService.VerifyEmail(c, req.Email, req.Code)
	if err != nil {
		status := http.StatusInternalServerError
		message := "验证邮箱失败"

		switch err {
		case service.ErrEmailNotFound:
			status = http.StatusNotFound
			message = "邮箱不存在"
		case service.ErrInvalidVerificationCode:
			status = http.StatusBadRequest
			message = "验证码无效或已过期"
		}

		c.JSON(status, dto.NewErrorResponse(status, message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.EmptySuccessResponse("邮箱验证成功"))
}

// ResendVerificationCode 重新发送验证码
func (h *EmailHandler) ResendVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
		return
	}

	_, err := h.userService.ResendVerificationCode(c, req.Email)
	if err != nil {
		status := http.StatusInternalServerError
		message := "重新发送验证码失败"

		if err == service.ErrEmailNotFound {
			status = http.StatusNotFound
			message = "邮箱不存在"
		}

		c.JSON(status, dto.NewErrorResponse(status, message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.EmptySuccessResponse("验证码已重新发送"))
}

// ForgotPassword 处理忘记密码请求
func (h *EmailHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
		return
	}

	_, err := h.userService.ForgotPassword(c, req.Email)
	if err != nil {
		status := http.StatusInternalServerError
		message := "处理忘记密码请求失败"

		if err == service.ErrEmailNotFound {
			status = http.StatusNotFound
			message = "邮箱不存在"
		}

		c.JSON(status, dto.NewErrorResponse(status, message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.EmptySuccessResponse("密码重置邮件已发送"))
}

// ResetPassword 处理密码重置请求
func (h *EmailHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "请求格式错误", err.Error()))
		return
	}

	err := h.userService.ResetPassword(c, req.Email, req.Token, req.NewPassword)
	if err != nil {
		status := http.StatusInternalServerError
		message := "重置密码失败"

		switch err {
		case service.ErrEmailNotFound:
			status = http.StatusNotFound
			message = "邮箱不存在"
		case service.ErrInvalidResetToken:
			status = http.StatusBadRequest
			message = "重置令牌无效或已过期"
		}

		c.JSON(status, dto.NewErrorResponse(status, message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.EmptySuccessResponse("密码重置成功"))
}
