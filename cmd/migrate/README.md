# 数据库迁移工具

这是一个独立的数据库迁移工具，用于管理 Memoir API 的数据库表结构。

## 功能特性

- **独立运行**：不依赖 API 服务器，可以单独执行
- **多种操作**：支持迁移、回滚、状态检查
- **安全可靠**：使用 GORM 的 AutoMigrate 和 Migrator 功能
- **配置统一**：使用与 API 相同的配置文件
- **明确操作**：必须明确指定操作类型，避免意外执行

## 使用方法

### 基本命令

```bash
# 运行数据库迁移（创建/更新表）
go run cmd/migrate/main.go -action=up

# 回滚数据库迁移（删除表）
go run cmd/migrate/main.go -action=down

# 检查迁移状态
go run cmd/migrate/main.go -action=status

# 显示帮助信息
go run cmd/migrate/main.go -help
```

### 使用 Makefile（推荐）

```bash
# 运行迁移
make migrate-up

# 回滚迁移
make migrate-down

# 检查状态
make migrate-status
```

## 迁移操作说明

### 1. 迁移 (up)
- 创建所有必要的数据库表
- 添加新的列和索引
- 不会删除现有的表或列
- 适合开发和生产环境的数据库更新

### 2. 回滚 (down)
- 删除所有应用程序相关的表
- 按照依赖关系逆序删除
- **注意：这会删除所有数据，请谨慎使用**

### 3. 状态检查 (status)
- 检查所有表是否存在
- 显示当前数据库的迁移状态
- 不会修改任何数据

## 配置要求

迁移工具使用与 API 服务器相同的配置：

- 确保 `.env` 文件存在并配置正确
- 数据库连接信息必须有效
- 数据库用户需要有创建/删除表的权限

## 部署建议

### 开发环境
```bash
# 首次设置
make migrate-up

# 开发过程中
make migrate-status  # 检查状态
make migrate-up      # 应用新的迁移
```

### 生产环境
```bash
# 部署前备份数据库
# 然后运行迁移
make migrate-up

# 验证迁移结果
make migrate-status
```

## 注意事项

1. **明确操作**：必须明确指定 `-action` 参数，工具不会执行默认操作
2. **备份数据**：在生产环境运行迁移前，请务必备份数据库
3. **权限检查**：确保数据库用户有足够的权限执行 DDL 操作
4. **回滚风险**：`migrate-down` 会删除所有数据，请谨慎使用
5. **并发安全**：避免同时运行多个迁移进程
