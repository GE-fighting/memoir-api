package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"memoir-api/internal/cache"
)

// STSConfig holds the configuration for Aliyun STS token generation
type STSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	RoleArn         string
	SessionName     string
	RegionID        string
	BucketName      string
}

// STSToken represents the token response structure
type STSToken struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Expiration      string `json:"expiration"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

// GetSTSConfig loads STS configuration from environment variables
func GetSTSConfig() (*STSConfig, error) {
	log.Println("开始加载STS配置...")

	accessKeyID := os.Getenv("ALIYUN_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	roleArn := os.Getenv("ALIYUN_ROLE_ARN")
	regionID := os.Getenv("ALIYUN_REGION_ID")
	bucketName := os.Getenv("ALIYUN_BUCKET_NAME")

	log.Printf("环境变量检查: ALIYUN_ACCESS_KEY_ID=%s, ALIYUN_ROLE_ARN=%s, ALIYUN_REGION_ID=%s, ALIYUN_BUCKET_NAME=%s",
		maskString(accessKeyID), roleArn, regionID, bucketName)

	if accessKeyID == "" || accessKeySecret == "" || roleArn == "" || regionID == "" || bucketName == "" {
		log.Println("错误: 缺少必需的环境变量")
		return nil, fmt.Errorf("missing required environment variables for Aliyun STS")
	}

	log.Println("STS配置加载成功")
	return &STSConfig{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		RoleArn:         roleArn,
		SessionName:     "memoir-session", // Fixed session name
		RegionID:        regionID,
		BucketName:      bucketName,
	}, nil
}

// GenerateSTSToken generates a temporary STS token for OSS access
// ownerID is used to scope the token to a specific user's directory
func GenerateSTSToken(ctx context.Context, ownerID string) (*STSToken, error) {

	//1、从redis中获取sts token
	stsToken, err := GetSTSTokenFromRedis(ctx, ownerID)
	if err != nil {
		log.Printf("从Redis获取STS令牌失败: %v", err)
	}
	if stsToken != nil {
		log.Printf("从Redis获取STS令牌成功: %v", stsToken)
		return stsToken, nil
	}

	log.Printf("开始为用户 %s 生成STS令牌...", ownerID)

	config, err := GetSTSConfig()
	if err != nil {
		log.Printf("获取STS配置失败: %v", err)
		return nil, err
	}

	// 创建OpenAPI配置
	clientConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyID),
		AccessKeySecret: tea.String(config.AccessKeySecret),
	}

	// 使用正确的区域端点
	endpoint := fmt.Sprintf("sts.%s.aliyuncs.com", config.RegionID)
	log.Printf("使用STS端点: %s", endpoint)
	clientConfig.Endpoint = tea.String(endpoint)

	// 创建STS客户端
	client, err := sts20150401.NewClient(clientConfig)
	if err != nil {
		log.Printf("创建STS客户端失败: %v", err)
		return nil, fmt.Errorf("failed to create STS client: %w", err)
	}

	// 创建AssumeRole请求
	policyStr := generatePolicy(config.BucketName, ownerID)
	request := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(config.RoleArn),
		RoleSessionName: tea.String(config.SessionName),
		DurationSeconds: tea.Int64(3600), // 1小时
		Policy:          tea.String(policyStr),
	}

	// 发送请求
	log.Println("发送AssumeRole请求...")
	response, err := client.AssumeRoleWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		log.Printf("AssumeRole请求失败: %v", err)
		return nil, fmt.Errorf("failed to assume role: %w", err)
	}

	log.Println("AssumeRole请求成功，已获取临时凭证")
	credentials := response.Body.Credentials

	// 创建token响应
	token := &STSToken{
		AccessKeyID:     *credentials.AccessKeyId,
		AccessKeySecret: *credentials.AccessKeySecret,
		SecurityToken:   *credentials.SecurityToken,
		Expiration:      *credentials.Expiration,
		Region:          config.RegionID,
		Bucket:          config.BucketName,
	}

	log.Println("STS令牌生成成功，过期时间:", *credentials.Expiration)

	//2、将sts token存入redis中
	err = SetSTSTokenToRedis(ctx, ownerID, token)
	if err != nil {
		log.Printf("将STS令牌存入Redis失败: %v", err)
	}
	return token, nil
}

// generatePolicy creates a policy document that restricts access to a specific user's path
func generatePolicy(bucket, ownerID string) string {
	log.Printf("为 %s 在存储桶 %s 生成策略...", ownerID, bucket)

	// 直接使用字符串模板创建策略JSON
	// 替换examplebucket/src/*为实际的bucket/userID/*

	resourcePath := fmt.Sprintf("acs:oss:*:*:%s/%s/*", bucket, ownerID)
	log.Printf("资源路径: %s", resourcePath)

	policyStr := fmt.Sprintf(`{
    "Version": "1", 
    "Statement": [
        {
            "Action": [
                "oss:PutObject",
				"oss:GetObject",
				"oss:DeleteObject"
            ], 
            "Resource": [
                "%s"
            ], 
            "Effect": "Allow"
        }
    ]
}`, resourcePath)

	log.Printf("完整策略JSON: %s", policyStr)
	return policyStr
}

// 工具函数，用于屏蔽敏感信息
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

func GetSTSTokenFromRedis(ctx context.Context, ownerID string) (*STSToken, error) {
	//TODO 从redis中获取sts token
	stsTokenStr, err := cache.GetRedisClient().Get(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	stsToken := &STSToken{}
	err = json.Unmarshal([]byte(stsTokenStr), stsToken)
	if err != nil {
		return nil, err
	}
	return stsToken, nil
}

func SetSTSTokenToRedis(ctx context.Context, ownerID string, stsToken *STSToken) error {
	//TODO 将sts token存入redis中
	return cache.GetRedisClient().Set(ctx, ownerID, stsToken, 3600*time.Second)
}
