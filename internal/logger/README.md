# Logger 使用指南

本包提供了基于 zerolog 的结构化日志系统，支持依赖注入和上下文感知功能。

## 功能特性

- **结构化日志**: 使用 JSON 格式输出日志，便于解析和分析
- **上下文感知**: 支持在请求上下文中传递日志实例，自动包含请求 ID 等关联信息
- **依赖注入**: 可以将日志实例注入到服务和仓库层
- **性能优化**: 基于高性能的 zerolog 库，内存分配最小化
- **分级日志**: 支持 Debug、Info、Warn、Error、Fatal 多个日志级别
- **丰富的上下文**: 支持添加字段、组件名称和错误信息

## 基本用法

### 初始化日志系统

在应用启动时，需要初始化日志系统：

```go
// 在 main.go 中初始化
cfg := config.New()
logger.Initialize(cfg.Server.LogLevel)
```

### 基本日志记录

```go
// 直接使用包级函数
logger.Info("服务启动成功", map[string]interface{}{
    "port": 8080,
})

// 记录错误
if err != nil {
    logger.Error(err, "数据库连接失败")
}
```

### 获取组件专用日志记录器

```go
// 为特定组件创建日志记录器
log := logger.GetLogger("user_service")

// 使用组件日志记录器
log.Info("用户登录成功", map[string]interface{}{
    "user_id": user.ID,
})
```

## 高级用法

### 上下文感知日志

在 HTTP 处理流程中，通常需要在整个请求生命周期内追踪相同的请求 ID。

#### 在中间件中设置上下文日志

```go
// 在 LoggerMiddleware 中
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取或生成请求 ID
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // 创建带请求 ID 的日志记录器
        reqLogger := logger.GetLogger("http").With("request_id", requestID)
        
        // 将日志记录器存入上下文
        c.Request = c.Request.WithContext(reqLogger.WithContext(c.Request.Context()))
        
        // 处理请求...
        c.Next()
    }
}
```

#### 在服务或仓库层中使用上下文日志

```go
func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
    // 从上下文中获取日志记录器
    log := logger.FromContext(ctx).WithComponent("user_service")
    
    log.Debug("获取用户信息", map[string]interface{}{
        "user_id": id,
    })
    
    // 处理业务逻辑...
    
    return user, nil
}
```

### 添加字段和上下文

```go
// 添加单个字段
log := logger.GetLogger("auth").With("module", "oauth2")

// 添加多个字段
log = log.WithFields(map[string]interface{}{
    "client_id": clientID,
    "scope": scope,
})

// 添加错误信息
if err != nil {
    log.WithError(err).Error("认证失败")
}
```

### 在多层应用中传递日志上下文

```go
// 控制器层
func (h *UserHandler) GetUser(c *gin.Context) {
    ctx := c.Request.Context()
    userID := getUserIDFromParams(c)
    
    user, err := h.userService.GetUser(ctx, userID)
    // ...
}

// 服务层
func (s *userService) GetUser(ctx context.Context, userID int64) (*models.User, error) {
    log := logger.FromContext(ctx).WithComponent("user_service")
    log.Debug("获取用户", map[string]interface{}{"user_id": userID})
    
    // 传递相同的上下文到仓库层
    return s.userRepo.GetByID(ctx, userID)
}

// 仓库层
func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    // 从相同的上下文中获取日志记录器
    log := logger.FromContext(ctx).WithComponent("user_repository")
    // ...
}
```

## 最佳实践

1. **总是使用上下文**: 在服务和仓库层方法中，总是接受 `context.Context` 参数并传递给下游函数
2. **组件标记**: 使用 `WithComponent()` 方法标记日志来源
3. **结构化字段**: 使用结构化字段而不是字符串拼接
4. **适当的日志级别**:
   - DEBUG: 详细的开发信息，用于调试
   - INFO: 正常的操作信息
   - WARN: 可能的问题但不影响基本功能
   - ERROR: 错误但应用可以继续运行
   - FATAL: 严重错误导致应用终止
5. **包含上下文**: 记录足够的上下文信息以便排查问题，如用户ID、请求ID等

## 配置

在 `.env` 文件中配置日志级别:

```
LOG_LEVEL=debug  # 可选: debug, info, warn, error, fatal
```

## 测试

提供了 `SetOutput` 函数用于测试:

```go
func TestSomething(t *testing.T) {
    // 捕获日志输出
    buf := &bytes.Buffer{}
    logger.SetOutput(buf)
    
    // 执行测试...
    
    // 验证日志输出
    logOutput := buf.String()
    if !strings.Contains(logOutput, "expected message") {
        t.Error("Expected log message not found")
    }
}
``` 