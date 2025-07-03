package email

import (
	"context"
	"encoding/json"
	"fmt"
	"memoir-api/internal/config"
	"memoir-api/internal/logger"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-redis/redis/v8"
)

// 邮件任务结构
type EmailTask struct {
	Type       EmailType         `json:"type"`
	ToAddress  string            `json:"to_address"`
	Subject    string            `json:"subject"`
	HtmlBody   string            `json:"html_body"`
	TextBody   string            `json:"text_body"`
	Data       map[string]string `json:"data"`
	RetryCount int               `json:"retry_count"`
	CreatedAt  time.Time         `json:"created_at"`
}

// EmailType 邮件类型
type EmailType string

// 邮件类型常量
const (
	EmailTypeVerification  EmailType = "verification"   // 邮箱验证
	EmailTypeResetPassword EmailType = "reset_password" // 密码重置
	EmailTypeNotification  EmailType = "notification"   // 系统通知
	EmailTypeWelcome       EmailType = "welcome"        // 欢迎邮件
	EmailTypeAnniversary   EmailType = "anniversary"    // 纪念日邮件
	EmailTypeFestival      EmailType = "festival"       // 节日邮件
)

// EmailQueue Redis队列名
const EmailQueue = "email:queue"

// RateLimit 键前缀
const RateLimitPrefix = "email:rate_limit:"

// 验证码有效期（分钟）
const VerificationCodeExpiry = 15

// EmailService 邮件服务接口
type EmailService interface {
	// 发送验证邮件
	SendVerificationEmail(ctx context.Context, toAddress, username, verificationCode string) error

	// 发送密码重置邮件
	SendPasswordResetEmail(ctx context.Context, toAddress, resetToken string) error

	// 发送通知邮件
	SendNotificationEmail(ctx context.Context, toAddress, username, message string) error

	// 发送欢迎邮件
	SendWelcomeEmail(ctx context.Context, toAddress, username string) error

	// 发送纪念日邮件
	SendAnniversaryEmail(ctx context.Context, toAddress, username, partnerName string, days int, date string) error

	// 发送节日邮件
	SendFestivalEmail(ctx context.Context, toAddress, username, partnerName, festivalName string) error

	// 处理邮件队列
	ProcessEmailQueue(ctx context.Context)

	// 存储验证码到Redis
	StoreVerificationCode(ctx context.Context, email, code string) error

	// 验证验证码
	VerifyCode(ctx context.Context, email, code string) (bool, error)

	// 存储密码重置令牌
	StorePasswordResetToken(ctx context.Context, email, token string) error

	// 验证密码重置令牌
	VerifyPasswordResetToken(ctx context.Context, email, token string) (bool, error)
}

// DirectMailService 阿里云邮件服务实现
type DirectMailService struct {
	client *dm20151123.Client
	config *config.EmailConfig
	redis  *redis.Client
	log    logger.Logger
}

// NewEmailService 创建新的邮件服务实例
func NewEmailService(cfg *config.Config, redisClient *redis.Client) (EmailService, error) {
	// 如果邮件服务未启用，返回空实现
	if !cfg.Email.Enabled {
		return &noOpEmailService{}, nil
	}

	// 创建阿里云邮件客户端
	client, err := createDMClient(cfg.Email)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云邮件客户端失败: %w", err)
	}

	return &DirectMailService{
		client: client,
		config: &cfg.Email,
		redis:  redisClient,
		log:    logger.GetLogger("email-service"),
	}, nil
}

// 创建阿里云邮件客户端
func createDMClient(cfg config.EmailConfig) (*dm20151123.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyID),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		RegionId:        tea.String(cfg.RegionID),
	}

	// 使用邮件服务的默认域名
	config.Endpoint = tea.String("dm.aliyuncs.com")

	return dm20151123.NewClient(config)
}

// SendVerificationEmail 发送验证邮件
func (s *DirectMailService) SendVerificationEmail(ctx context.Context, toAddress, username, verificationCode string) error {
	// 检查发送频率限制
	if err := s.checkRateLimit(ctx, toAddress); err != nil {
		return err
	}

	// 存储验证码到Redis
	if err := s.StoreVerificationCode(ctx, toAddress, verificationCode); err != nil {
		s.log.Error(err, "存储验证码失败")
		return err
	}

	// 准备邮件内容
	task := EmailTask{
		Type:      EmailTypeVerification,
		ToAddress: toAddress,
		Subject:   fmt.Sprintf("%s - 请验证您的邮箱", s.config.AppName),
		Data: map[string]string{
			"AppName":          s.config.AppName,
			"Username":         username,
			"VerificationCode": verificationCode,
			"ExpireMinutes":    fmt.Sprintf("%d", VerificationCodeExpiry),
		},
		CreatedAt: time.Now(),
	}

	// 渲染邮件内容
	task.HtmlBody = renderVerificationEmailTemplate(task.Data)
	task.TextBody = fmt.Sprintf("您好 %s，您的验证码是：%s，%d分钟内有效。",
		username, verificationCode, VerificationCodeExpiry)

	return s.addToQueue(ctx, task)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *DirectMailService) SendPasswordResetEmail(ctx context.Context, toAddress, resetToken string) error {
	// 检查发送频率限制
	if err := s.checkRateLimit(ctx, toAddress); err != nil {
		return err
	}

	// 存储密码重置令牌到Redis
	if err := s.StorePasswordResetToken(ctx, toAddress, resetToken); err != nil {
		s.log.Error(err, "存储密码重置令牌失败")
		return err
	}

	// 构建重置链接
	resetLink := fmt.Sprintf("%s/reset-password?token=%s&email=%s",
		s.config.AppURL, resetToken, toAddress)

	// 准备邮件内容
	task := EmailTask{
		Type:      EmailTypeResetPassword,
		ToAddress: toAddress,
		Subject:   fmt.Sprintf("%s - 密码重置请求", s.config.AppName),
		Data: map[string]string{
			"AppName":   s.config.AppName,
			"ResetLink": resetLink,
		},
		CreatedAt: time.Now(),
	}

	// 渲染邮件内容
	task.HtmlBody = renderPasswordResetEmailTemplate(task.Data)
	task.TextBody = fmt.Sprintf("您请求重置密码，请访问以下链接完成重置：%s", resetLink)

	return s.addToQueue(ctx, task)
}

// SendNotificationEmail 发送通知邮件
func (s *DirectMailService) SendNotificationEmail(ctx context.Context, toAddress, username, message string) error {
	// 检查发送频率限制
	if err := s.checkRateLimit(ctx, toAddress); err != nil {
		return err
	}

	// 准备邮件内容
	task := EmailTask{
		Type:      EmailTypeNotification,
		ToAddress: toAddress,
		Subject:   fmt.Sprintf("%s - 系统通知", s.config.AppName),
		Data: map[string]string{
			"AppName":  s.config.AppName,
			"Username": username,
			"Message":  message,
		},
		CreatedAt: time.Now(),
	}

	// 渲染邮件内容
	task.HtmlBody = renderNotificationEmailTemplate(task.Data)
	task.TextBody = fmt.Sprintf("您好 %s，%s", username, message)

	return s.addToQueue(ctx, task)
}

// SendWelcomeEmail 发送欢迎邮件
func (s *DirectMailService) SendWelcomeEmail(ctx context.Context, toAddress, username string) error {
	// 准备邮件内容
	task := EmailTask{
		Type:      EmailTypeWelcome,
		ToAddress: toAddress,
		Subject:   fmt.Sprintf("欢迎加入 %s", s.config.AppName),
		Data: map[string]string{
			"AppName":  s.config.AppName,
			"Username": username,
			"AppURL":   s.config.AppURL,
		},
		CreatedAt: time.Now(),
	}

	// 渲染邮件内容
	task.HtmlBody = renderWelcomeEmailTemplate(task.Data)
	task.TextBody = fmt.Sprintf("欢迎 %s 加入 %s！", username, s.config.AppName)

	return s.addToQueue(ctx, task)
}

// SendAnniversaryEmail 发送纪念日邮件
func (s *DirectMailService) SendAnniversaryEmail(ctx context.Context, toAddress, username, partnerName string, days int, date string) error {
	// 准备邮件内容
	task := EmailTask{
		Type:      EmailTypeAnniversary,
		ToAddress: toAddress,
		Subject:   fmt.Sprintf("❤️ 您与%s的恋爱纪念日 - %s天快乐", partnerName, formatDays(days)),
		Data: map[string]string{
			"AppName":     s.config.AppName,
			"Username":    username,
			"PartnerName": partnerName,
			"Days":        formatDays(days),
			"Date":        date,
			"AppURL":      s.config.AppURL,
		},
		CreatedAt: time.Now(),
	}

	// 渲染邮件内容
	task.HtmlBody = renderAnniversaryEmailTemplate(task.Data)
	task.TextBody = fmt.Sprintf("亲爱的%s，今天是您和%s在一起的第%s天！祝福你们爱情长久，记得珍惜这美好的时光。",
		username, partnerName, formatDays(days))

	return s.addToQueue(ctx, task)
}

// SendFestivalEmail 发送节日邮件
func (s *DirectMailService) SendFestivalEmail(ctx context.Context, toAddress, username, partnerName, festivalName string) error {
	// 准备邮件内容
	task := EmailTask{
		Type:      EmailTypeFestival,
		ToAddress: toAddress,
		Subject:   fmt.Sprintf("❤️ %s快乐 - 给%s的祝福", festivalName, username),
		Data: map[string]string{
			"AppName":      s.config.AppName,
			"Username":     username,
			"PartnerName":  partnerName,
			"FestivalName": festivalName,
			"AppURL":       s.config.AppURL,
		},
		CreatedAt: time.Now(),
	}

	// 渲染邮件内容
	task.HtmlBody = renderFestivalEmailTemplate(task.Data)
	task.TextBody = fmt.Sprintf("亲爱的%s，祝您和%s%s快乐！希望你们能一起度过一个浪漫美好的节日。",
		username, partnerName, festivalName)

	return s.addToQueue(ctx, task)
}

// 格式化天数，添加特殊处理
func formatDays(days int) string {
	if days == 100 {
		return "100"
	} else if days == 365 || days == 366 {
		return "一周年"
	} else if days == 730 || days == 731 {
		return "两周年"
	} else if days == 1095 || days == 1096 {
		return "三周年"
	} else if days == 1825 || days == 1826 {
		return "五周年"
	} else if days == 3650 || days == 3651 || days == 3652 {
		return "十周年"
	} else {
		return fmt.Sprintf("%d", days)
	}
}

// 添加任务到队列
func (s *DirectMailService) addToQueue(ctx context.Context, task EmailTask) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		s.log.Error(err, "序列化邮件任务失败")
		return fmt.Errorf("序列化邮件任务失败: %w", err)
	}

	err = s.redis.LPush(ctx, EmailQueue, taskJSON).Err()
	if err != nil {
		s.log.Error(err, "添加邮件任务到队列失败")
		return fmt.Errorf("添加邮件任务到队列失败: %w", err)
	}

	s.log.Info("邮件任务已添加到队列", "type", task.Type, "to", task.ToAddress)
	return nil
}

// ProcessEmailQueue 处理邮件队列
func (s *DirectMailService) ProcessEmailQueue(ctx context.Context) {
	s.log.Info("开始处理邮件队列")

	for {
		select {
		case <-ctx.Done():
			s.log.Info("邮件队列处理已停止")
			return
		default:
			// 从队列获取任务
			result, err := s.redis.BRPop(ctx, 5*time.Second, EmailQueue).Result()
			if err != nil {
				if err != redis.Nil {
					s.log.Error(err, "从队列获取任务失败")
				}
				continue
			}

			if len(result) < 2 {
				continue
			}

			var task EmailTask
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				s.log.Error(err, "反序列化邮件任务失败")
				continue
			}

			// 发送邮件
			s.log.Info("处理邮件任务", "type", task.Type, "to", task.ToAddress)
			if err := s.sendEmail(task); err != nil {
				s.log.Error(err, "发送邮件失败", "type", task.Type, "to", task.ToAddress)

				// 重试逻辑
				if task.RetryCount < 3 {
					task.RetryCount++
					s.log.Info("邮件发送失败，稍后重试", "retry", task.RetryCount, "to", task.ToAddress)
					time.Sleep(time.Duration(task.RetryCount) * time.Minute)
					s.addToQueue(ctx, task)
				} else {
					s.log.Error(nil, "邮件发送失败，超过最大重试次数", "to", task.ToAddress)
				}
			}
		}
	}
}

// 实际发送邮件
func (s *DirectMailService) sendEmail(task EmailTask) error {
	s.log.Info("发送邮件", "to", task.ToAddress, "subject", task.Subject)

	request := &dm20151123.SingleSendMailRequest{
		AccountName:    tea.String(s.config.AccountName),
		AddressType:    tea.Int32(int32(s.config.AddressType)),
		ReplyToAddress: tea.Bool(s.config.ReplyToAddress),
		ToAddress:      tea.String(task.ToAddress),
		Subject:        tea.String(task.Subject),
		HtmlBody:       tea.String(task.HtmlBody),
		TextBody:       tea.String(task.TextBody),
		FromAlias:      tea.String(s.config.FromAlias),
	}

	response, err := s.client.SingleSendMail(request)
	if err != nil {
		return fmt.Errorf("调用阿里云邮件API失败: %w", err)
	}

	if response.StatusCode == nil || *response.StatusCode != 200 {
		// 提取错误信息
		errMsg := "未知错误"

		// 如果有响应体和RequestId，则加入错误信息
		if response.Body != nil && response.Body.RequestId != nil {
			errMsg = fmt.Sprintf("请求失败 (RequestId: %s)", *response.Body.RequestId)
		}
		return fmt.Errorf("邮件发送失败: %s", errMsg)
	}

	s.log.Info("邮件发送成功", "type", task.Type, "to", task.ToAddress, "requestId", *response.Body.RequestId)
	return nil
}

// 检查发送频率限制
func (s *DirectMailService) checkRateLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("%s%s", RateLimitPrefix, email)
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		s.log.Error(err, "检查发送频率失败")
		return err
	}

	if count == 1 {
		s.redis.Expire(ctx, key, time.Hour)
	}

	if count > 5 { // 每小时最多5封邮件
		return fmt.Errorf("发送频率过高，请稍后再试")
	}

	return nil
}

// StoreVerificationCode 存储验证码到Redis
func (s *DirectMailService) StoreVerificationCode(ctx context.Context, email, code string) error {
	key := fmt.Sprintf("email:verify:%s", email)
	return s.redis.Set(ctx, key, code, VerificationCodeExpiry*time.Minute).Err()
}

// VerifyCode 验证验证码
func (s *DirectMailService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	key := fmt.Sprintf("email:verify:%s", email)
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if storedCode != code {
		return false, nil
	}

	// 验证成功后删除验证码
	s.redis.Del(ctx, key)
	return true, nil
}

// StorePasswordResetToken 存储密码重置令牌
func (s *DirectMailService) StorePasswordResetToken(ctx context.Context, email, token string) error {
	key := fmt.Sprintf("email:reset:%s", email)
	return s.redis.Set(ctx, key, token, 30*time.Minute).Err()
}

// VerifyPasswordResetToken 验证密码重置令牌
func (s *DirectMailService) VerifyPasswordResetToken(ctx context.Context, email, token string) (bool, error) {
	key := fmt.Sprintf("email:reset:%s", email)
	storedToken, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if storedToken != token {
		return false, nil
	}

	// 验证成功后删除令牌（一次性使用）
	s.redis.Del(ctx, key)
	return true, nil
}

// noOpEmailService 空实现（当邮件服务未启用时使用）
type noOpEmailService struct{}

func (s *noOpEmailService) SendVerificationEmail(ctx context.Context, toAddress, username, verificationCode string) error {
	return nil
}

func (s *noOpEmailService) SendPasswordResetEmail(ctx context.Context, toAddress, resetToken string) error {
	return nil
}

func (s *noOpEmailService) SendNotificationEmail(ctx context.Context, toAddress, username, message string) error {
	return nil
}

func (s *noOpEmailService) SendWelcomeEmail(ctx context.Context, toAddress, username string) error {
	return nil
}

func (s *noOpEmailService) SendAnniversaryEmail(ctx context.Context, toAddress, username, partnerName string, days int, date string) error {
	return nil
}

func (s *noOpEmailService) SendFestivalEmail(ctx context.Context, toAddress, username, partnerName, festivalName string) error {
	return nil
}

func (s *noOpEmailService) ProcessEmailQueue(ctx context.Context) {
	// 空实现，不做任何处理
}

func (s *noOpEmailService) StoreVerificationCode(ctx context.Context, email, code string) error {
	return nil
}

func (s *noOpEmailService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	return false, nil
}

func (s *noOpEmailService) StorePasswordResetToken(ctx context.Context, email, token string) error {
	return nil
}

func (s *noOpEmailService) VerifyPasswordResetToken(ctx context.Context, email, token string) (bool, error) {
	return false, nil
}
