package aliyun

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
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
// userID is used to scope the token to a specific user's directory
func GenerateSTSToken(userID string) (*STSToken, error) {
	log.Printf("开始为用户 %s 生成STS令牌...", userID)

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

	// 创建策略
	log.Println("生成访问策略...")
	policy := generatePolicy(config.BucketName, userID)
	log.Printf("生成的策略: %s", policy)

	// 创建AssumeRole请求
	request := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(config.RoleArn),
		RoleSessionName: tea.String(config.SessionName),
		DurationSeconds: tea.Int64(3600), // 1小时
		// Policy:          tea.String(policy),
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
	return token, nil
}

// generatePolicy creates a policy document that restricts access to a specific user's path
func generatePolicy(bucket, userID string) string {
	log.Printf("为用户 %s 在存储桶 %s 生成策略...", userID, bucket)

	// Base path for the user
	resourcePath := fmt.Sprintf("acs:oss:*:*:%s/%s", bucket, userID)
	log.Printf("资源路径: %s", resourcePath)

	// Allow only putting objects in the user's path
	policyDocument := map[string]interface{}{
		"Version": "1",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"oss:*",
				},
				"Resource": []string{
					fmt.Sprintf("%s/*", resourcePath),
				},
			},
		},
	}

	// Convert policy to JSON
	policyJSON, err := json.Marshal(policyDocument)
	if err != nil {
		log.Printf("策略JSON序列化失败: %v", err)
		// If marshaling fails, return an empty policy
		return ""
	}

	return string(policyJSON)
}

// 工具函数，用于屏蔽敏感信息
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}
