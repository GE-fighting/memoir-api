package cache

import (
	"log"
	"memoir-api/internal/config"
	"sync"
)

var (
	redisClient *Redis
	once        sync.Once
)

// Init 初始化Redis连接
func Init(cfg *config.Config) error {
	var err error

	once.Do(func() {
		redisCfg := &RedisConfig{
			Host:     cfg.Redis.Host,
			Port:     cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}
		log.Printf("Redis配置: %+v", redisCfg)
		redisClient, err = NewRedis(redisCfg)
		if err != nil {
			log.Printf("Redis初始化失败: %v", err)
			return
		}

		log.Printf("Redis连接成功: %s:%s", cfg.Redis.Host, cfg.Redis.Port)
	})

	return err
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *Redis {
	if redisClient == nil {
		log.Println("警告: Redis客户端尚未初始化")
	}
	return redisClient
}

// Close 关闭Redis连接
func Close() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
