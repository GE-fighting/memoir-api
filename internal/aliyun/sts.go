package aliyun

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
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

	// Create STS client
	log.Println("创建STS客户端...")
	client, err := sts.NewClientWithAccessKey(
		config.RegionID,
		config.AccessKeyID,
		config.AccessKeySecret,
	)
	if err != nil {
		log.Printf("创建STS客户端失败: %v", err)
		return nil, fmt.Errorf("failed to create STS client: %w", err)
	}

	// Create request
	log.Println("创建AssumeRole请求...")
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = config.RoleArn
	request.RoleSessionName = config.SessionName

	// Set token expiration to 1 hour
	request.DurationSeconds = "3600"

	// Create policy to restrict access to user's directory
	log.Println("生成访问策略...")
	policy := generatePolicy(config.BucketName, userID)
	log.Printf("生成的策略: %s", policy)
	request.Policy = policy

	// Send request
	log.Println("发送AssumeRole请求...")
	response, err := client.AssumeRole(request)
	if err != nil {
		log.Printf("AssumeRole请求失败: %v", err)
		return nil, fmt.Errorf("failed to assume role: %w", err)
	}

	log.Println("AssumeRole请求成功，已获取临时凭证")

	// Create token response
	token := &STSToken{
		AccessKeyID:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SecurityToken:   response.Credentials.SecurityToken,
		Expiration:      response.Credentials.Expiration,
		Region:          config.RegionID,
		Bucket:          config.BucketName,
	}

	log.Println("STS令牌生成成功，过期时间:", response.Credentials.Expiration)
	return token, nil
}

// generatePolicy creates a policy document that restricts access to a specific user's path
func generatePolicy(bucket, userID string) string {
	log.Printf("为用户 %s 在存储桶 %s 生成策略...", userID, bucket)

	// Base path for the user
	resourcePath := fmt.Sprintf("acs:oss:*:*:%s/%s", bucket, userID)
	log.Printf("资源路径: %s", resourcePath)

	// Allow listing objects, getting objects, and putting objects
	policyDocument := map[string]interface{}{
		"Version": "1",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"oss:ListObjects",
					"oss:GetObject",
					"oss:PutObject",
					"oss:DeleteObject",
				},
				"Resource": []string{
					resourcePath,
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
