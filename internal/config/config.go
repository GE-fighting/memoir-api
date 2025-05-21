package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 存储应用程序配置
type Config struct {
	DB    DBConfig
	Redis RedisConfig
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
	}
}

// 从环境变量获取值，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
