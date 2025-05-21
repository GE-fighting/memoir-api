package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig 存储Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Redis 客户端
type Redis struct {
	Client *redis.Client
}

// NewRedis 创建一个新的Redis客户端
func NewRedis(cfg *RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("无法连接到Redis: %v", err)
	}

	return &Redis{Client: client}, nil
}

// Set 存储键值对，可选过期时间
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var dataToStore string

	// 如果是字符串类型直接存储
	if str, ok := value.(string); ok {
		dataToStore = str
	} else {
		// 否则序列化为JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("序列化数据失败: %v", err)
		}
		dataToStore = string(jsonBytes)
	}

	return r.Client.Set(ctx, key, dataToStore, expiration).Err()
}

// Get 获取键对应的值
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// GetObject 获取键对应的值并反序列化为指定类型
func (r *Redis) GetObject(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除一个或多个键
func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetWithTTL 存储带有过期时间的键值对的简便方法
func (r *Redis) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Set(ctx, key, value, ttl)
}

// Close 关闭Redis连接
func (r *Redis) Close() error {
	return r.Client.Close()
}
