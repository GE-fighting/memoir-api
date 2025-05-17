package handlers

import (
	"memoir-api/internal/repository"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
)

// 获取用户ID从上下文中
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, ErrUnauthorized
	}

	return userID.(int64), nil
}

// 身份认证中间件
func AuthMiddleware() gin.HandlerFunc {
	// 使用已经定义的JWTAuthMiddleware
	return JWTAuthMiddleware()
}

// RegisterAllHandlers 注册所有处理器
func RegisterAllHandlers(router *gin.Engine, repoFactory repository.Factory, serviceFactory service.Factory) {
	// 注册个人媒体处理器
	RegisterPersonalMediaHandlers(router, repoFactory, serviceFactory)
}

// RegisterPersonalMediaHandlers 注册个人媒体处理器
func RegisterPersonalMediaHandlers(router *gin.Engine, repoFactory repository.Factory, serviceFactory service.Factory) {
	// 创建个人媒体服务
	personalMediaService := serviceFactory.PersonalMedia()

	// 创建个人媒体处理器
	personalMediaHandler := NewPersonalMediaHandler(personalMediaService)

	// 注册路由
	personalMediaHandler.RegisterRoutes(router)
}
