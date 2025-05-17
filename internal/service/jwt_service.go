package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 定义令牌类型
type TokenType string

const (
	// AccessToken 访问令牌
	AccessToken TokenType = "access"
	// RefreshToken 刷新令牌
	RefreshToken TokenType = "refresh"
)

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

// TokenDetails 令牌详情
type TokenDetails struct {
	AccessToken   string
	RefreshToken  string
	AccessExpiry  int64
	RefreshExpiry int64
	AccessUUID    string
	RefreshUUID   string
	ExpiresIn     int64
}

// JWTService JWT服务接口
type JWTService interface {
	Service
	// GenerateTokens 生成访问令牌和刷新令牌
	GenerateTokens(userID int64) (*TokenDetails, error)
	// ValidateToken 验证令牌
	ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	// ExtractUserID 从令牌中提取用户ID
	ExtractUserID(tokenString string) (int64, error)
}

// jwtService JWT服务实现
type jwtService struct {
	*BaseService
	config JWTConfig
}

// NewJWTService 创建JWT服务
func NewJWTService() JWTService {
	// 从环境变量获取密钥
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		// 默认开发密钥
		secretKey = "memoir-api-development-secret-key"
	}

	return &jwtService{
		BaseService: NewBaseService(nil),
		config: JWTConfig{
			SecretKey:     secretKey,
			AccessExpiry:  time.Hour * 24,     // 访问令牌有效期24小时
			RefreshExpiry: time.Hour * 24 * 7, // 刷新令牌有效期7天
			Issuer:        "memoir-api",
		},
	}
}

// GenerateTokens 生成访问令牌和刷新令牌
func (s *jwtService) GenerateTokens(userID int64) (*TokenDetails, error) {
	td := &TokenDetails{}
	now := time.Now()

	// 设置过期时间
	accessExpiry := now.Add(s.config.AccessExpiry)
	refreshExpiry := now.Add(s.config.RefreshExpiry)

	td.AccessExpiry = accessExpiry.Unix()
	td.RefreshExpiry = refreshExpiry.Unix()
	td.ExpiresIn = accessExpiry.Unix() - now.Unix()

	// 创建访问令牌
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": td.AccessExpiry,
		"iat": now.Unix(),
		"iss": s.config.Issuer,
		"typ": string(AccessToken),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	var err error
	td.AccessToken, err = accessToken.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("创建访问令牌失败: %w", err)
	}

	// 创建刷新令牌
	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": td.RefreshExpiry,
		"iat": now.Unix(),
		"iss": s.config.Issuer,
		"typ": string(RefreshToken),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("创建刷新令牌失败: %w", err)
	}

	return td, nil
}

// ValidateToken 验证令牌
func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("无效的令牌")
	}

	return token, claims, nil
}

// ExtractUserID 从令牌中提取用户ID
func (s *jwtService) ExtractUserID(tokenString string) (int64, error) {
	_, claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("无法获取用户ID")
	}

	return int64(userID), nil
}
