# Memoir API

情侣纪念相册后端服务，基于 Go 实现。

## 项目结构

```
memoir-api/
├── cmd/                      # 应用程序入口点
│   └── app/                  # 主应用
│       └── main.go           # 应用程序主入口
├── config/                   # 配置文件和配置加载
├── internal/                 # 私有应用程序代码
│   ├── api/                  # API 层
│   │   ├── dto/              # 数据传输对象
│   │   ├── errors/           # 错误定义和处理
│   │   ├── middleware/       # 中间件
│   │   └── routes.go         # 路由注册
│   ├── handlers/             # HTTP 处理器
│   ├── models/               # 数据模型
│   ├── repository/           # 数据访问层
│   ├── service/              # 业务逻辑层
│   └── utils/                # 工具函数
├── migrations/               # 数据库迁移文件
├── scripts/                  # 构建和部署脚本
└── tests/                    # 测试代码
```

## 核心组件

### 1. 数据传输对象 (DTO)

位于 `internal/api/dto` 目录，包含请求和响应的数据结构定义。

### 2. 错误处理

位于 `internal/api/errors` 目录，包含标准化的错误定义和处理函数。

### 3. 中间件

位于 `internal/api/middleware` 目录，包含:
- JWT 认证
- 请求日志
- CORS 配置
- 错误处理

### 4. 处理器

位于 `internal/handlers` 目录，处理 HTTP 请求并与服务层交互。

### 5. 路由注册

位于 `internal/api/routes.go`，将处理器与路由绑定。

## 环境变量

- `APP_ENV`: 应用环境 (development, production)
- `PORT`: 应用端口
- `DB_URL`: 数据库连接字符串
- `JWT_SECRET`: JWT 加密密钥

## 运行应用

```bash
go run cmd/app/main.go
```

## API 文档

访问 `/swagger/index.html` 查看 API 文档。

## 功能

- 时间轴记录：记录重要时刻，支持按年/月/日筛选
- 相册墙：支持照片和视频，多种布局展示
- 心愿清单：管理共同心愿，记录完成状态
- 纪念地图：记录去过的地方，支持地图标记

## 开发准备

### 配置数据库连接信息

1. 复制示例配置文件（如果尚未复制）：
   ```bash
   cp .env.example .env
   ```

2. 编辑 `.env` 文件，根据你的本地环境修改数据库连接信息：
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=你的密码
   DB_NAME=memoir
   DB_SSLMODE=disable
   LOG_LEVEL=debug   # 可选: debug, info, warn, error
   GIN_MODE=debug    # 可选: debug, release
   ```

3. 创建数据库（如果尚未创建）：
   ```sql
   CREATE DATABASE memoir;
   ```

> 注意：`.env` 文件不会提交到版本控制系统中，每个开发者需要维护自己的本地配置。

## 数据库管理

本项目使用 PostgreSQL 数据库存储数据。数据库脚本位于 `sql` 目录下。

### 前置条件

1. 确保 PostgreSQL 数据库已安装并运行
2. 创建名为 `memoir` 的数据库
3. 确保环境变量正确设置（见上方配置说明）

### 初始化数据库

执行以下命令创建数据库表结构：

```bash
psql -h <主机名> -p <端口> -U <用户名> -d memoir -f sql/schema.sql
```

示例：

```bash
psql -h localhost -p 5432 -U postgres -d memoir -f sql/schema.sql
```

### 清理数据库

如需清理所有表，可执行：

```bash
psql -h <主机名> -p <端口> -U <用户名> -d memoir -f sql/drop_tables.sql
```

### 自动迁移

本项目使用 GORM 的 `AutoMigrate` 功能进行自动数据库迁移，启动服务时会自动创建或更新表结构。这种方式适合开发环境，避免手动执行 SQL 脚本的麻烦。

## CORS 配置

为了解决跨域资源共享(CORS)问题，你需要在 `.env` 文件中配置允许的源。

创建一个 `.env` 文件在项目根目录，添加以下内容：

```
# CORS配置
# 多个源用逗号分隔
CORS_ORIGINS=http://localhost:3000,http://172.28.24.190:3000
```

你可以根据需要添加更多的源。

## 其他环境变量

完整的环境变量配置示例：

```
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=memoir
DB_SSLMODE=disable

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 服务器配置
SERVER_PORT=5000
SERVER_HOST=0.0.0.0
SERVER_MODE=debug  # debug 或 release
SERVER_READTIMEOUT=10
SERVER_WRITETIMEOUT=30
SERVER_IDLETIMEOUT=60
SERVER_LOGLEVEL=debug  # debug, info, warn, error
SERVER_MAXBODYSIZE=10485760  # 10MB
SERVER_JWTSECRET=your_jwt_secret_here
SERVER_JWTEXPIRE=24  # 小时

# CORS配置
CORS_ORIGINS=http://localhost:3000,http://172.28.24.190:3000
```

