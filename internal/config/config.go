package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 存储应用程序配置
type Config struct {
	DB     DBConfig
	Redis  RedisConfig
	Server ServerConfig
}

// DBConfig 存储数据库配置
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig 存储Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port         int      // 服务监听端口
	Host         string   // 绑定地址，如 "0.0.0.0" 或 "localhost"
	ReadTimeout  int      // HTTP读取超时时间(秒)
	WriteTimeout int      // HTTP写入超时时间(秒)
	IdleTimeout  int      // 保持连接超时时间(秒)
	Mode         string   // 运行模式：development, production
	LogLevel     string   // 日志级别
	MaxBodySize  int64    // 最大请求体大小(字节)
	JWTSecret    string   // JWT密钥(如果使用JWT认证)
	JWTExpire    int      // JWT过期时间(小时)
	CorsOrigins  []string // CORS允许的源
}

// ConnectionString 返回数据库连接字符串
func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// LoadEnv 加载.env文件中的环境变量
func LoadEnv() {
	// 获取项目根目录
	rootDir, err := findRootDir()
	if err != nil {
		log.Printf("无法获取项目根目录: %v", err)
		log.Println("尝试从当前目录加载.env文件")
		rootDir = "."
	}

	// 尝试加载.env文件（如果存在）
	envFile := filepath.Join(rootDir, ".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf(".env文件未找到或无法加载，将使用系统环境变量或默认值: %v", err)
	} else {
		log.Printf("已从 %s 加载环境变量配置", envFile)
	}
}

// 查找项目根目录（包含go.mod的目录）
func findRootDir() (string, error) {
	// 从当前目录开始
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上查找，直到找到包含go.mod的目录
	for {
		// 检查go.mod是否存在
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// 获取父目录
		parentDir := filepath.Dir(dir)
		// 如果已经到达根目录，则停止
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("无法找到项目根目录")
}

// New 返回应用程序配置
func New() *Config {
	// 首先尝试加载.env文件
	LoadEnv()

	// 解析Redis DB值
	redisDB := 0
	if dbStr := getEnv("REDIS_DB", "0"); dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			redisDB = parsed
		}
	}

	// 解析CORS Origins
	corsOrigins := []string{"http://localhost:3000"}
	if originsStr := getEnv("CORS_ORIGINS", ""); originsStr != "" {
		corsOrigins = append(corsOrigins, strings.Split(originsStr, ",")...)
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "memoir"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", "5000"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Mode:         getEnv("SERVER_MODE", "debug"),
			ReadTimeout:  getEnvInt("SERVER_READTIMEOUT", "10"),
			WriteTimeout: getEnvInt("SERVER_WRITETIMEOUT", "30"),
			IdleTimeout:  getEnvInt("SERVER_IDLETIMEOUT", "60"),
			LogLevel:     getEnv("SERVER_LOGLEVEL", "debug"),
			MaxBodySize:  getEnvInt64("SERVER_MAXBODYSIZE", "10485760"), // 默认10MB
			JWTSecret:    getEnv("SERVER_JWTSECRET", ""),
			JWTExpire:    getEnvInt("SERVER_JWTEXPIRE", "24"),
			CorsOrigins:  corsOrigins,
		},
	}
}

// 从环境变量获取值，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key, defaultValue string) int {
	valStr := getEnv(key, defaultValue)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0 // 或者返回你想要的默认值，比如 10
	}
	return val
}

func getEnvInt64(key, defaultValue string) int64 {
	valStr := getEnv(key, defaultValue)
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0 // 或者返回你想要的默认值，比如 10485760
	}
	return val
}
