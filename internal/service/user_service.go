package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// 错误常量定义
var (
	ErrUserExists              = errors.New("用户已存在")
	ErrInvalidPassword         = errors.New("密码不正确")
	ErrVerificationRequired    = errors.New("需要验证邮箱")
	ErrInvalidVerificationCode = errors.New("验证码无效")
	ErrInvalidResetToken       = errors.New("重置令牌无效")
	ErrEmailNotFound           = errors.New("邮箱不存在")
)

// UserService 用户服务接口
type UserService interface {
	Service
	Register(ctx context.Context, username, email, password, pairToken string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	LoginByUsername(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
	ExistCouple(ctx context.Context) (bool, error)
	GetCoupleID(ctx context.Context, userID int64) (int64, error)
	VerifyEmail(ctx context.Context, email, code string) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, email, token, newPassword string) error
	GenerateVerificationCode() string
	ResendVerificationCode(ctx context.Context, email string) (string, error)
}

// userService 用户服务实现
type userService struct {
	*BaseService
	userRepo   repository.UserRepository
	coupleRepo repository.CoupleRepository
	emailSvc   EmailService // 邮件服务依赖
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, coupleRepo repository.CoupleRepository, emailSvc EmailService) UserService {
	return &userService{
		BaseService: NewBaseService(userRepo),
		userRepo:    userRepo,
		coupleRepo:  coupleRepo,
		emailSvc:    emailSvc,
	}
}

// Register 注册用户
func (s *userService) Register(ctx context.Context, username, email, password, pairToken string) (*models.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	} else if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("检查用户邮箱时发生错误: %w", err)
	}

	// 检查用户名是否已存在
	existingUser, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	} else if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("检查用户名时发生错误: %w", err)
	}

	// 对密码进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}

	// 如果pair_token不为空，查询是否有情侣记录，如果有的话返回情侣主键id；没用的创建记录然后返回
	var coupleID int64
	if pairToken != "" {
		// 查找是否存在使用此配对令牌的情侣关系
		couple, err := s.coupleRepo.GetByPairToken(ctx, pairToken)
		if err != nil {
			if !errors.Is(err, repository.ErrCoupleNotFound) {
				return nil, fmt.Errorf("查询情侣记录时发生错误: %w", err)
			}
		}
		// 如果不存在,创建新的情侣关系
		if couple == nil {
			couple = &models.Couple{
				PairToken:             pairToken,
				AutoGenerateVideo:     true, // 默认开启
				ReminderNotifications: true, // 默认开启
			}
			if err := s.coupleRepo.Create(ctx, couple); err != nil {
				return nil, fmt.Errorf("创建情侣关系失败: %w", err)
			}
		}
		coupleID = couple.ID
	}

	// 创建用户
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CoupleID:     coupleID,
		DarkMode:     false, // 默认不开启
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 生成并发送验证码
	verificationCode := s.GenerateVerificationCode()
	err = s.emailSvc.SendVerificationEmail(ctx, email, username, verificationCode)
	if err != nil {
		// 记录错误但不影响注册流程
		fmt.Printf("发送验证邮件失败: %v", err)
	}

	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户时发生错误: %w", err)
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// LoginByUsername 通过用户名登录
func (s *userService) LoginByUsername(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户时发生错误: %w", err)
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// GetUserByID 通过ID获取用户
func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail 通过邮箱获取用户
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// UpdatePassword 更新用户密码
func (s *userService) UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// 对新密码进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

// GeneratePairToken 生成唯一的情侣配对令牌
func GeneratePairToken() string {
	// 生成一个 UUID 风格的令牌，12 位字母数字
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 12

	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果随机数生成失败，使用时间戳作为备选
			return time.Now().String()
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}

func (s *userService) ExistCouple(ctx context.Context) (bool, error) {
	userID, exists := ctx.Value("user_id").(int64)
	if !exists {
		return false, errors.New("user_id not found")
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.CoupleID != 0, nil
}

func (s *userService) GetCoupleID(ctx context.Context, userID int64) (int64, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if user.CoupleID == 0 {
		return 0, errors.New("用户没有情侣关系")
	}
	couple, err := s.coupleRepo.GetByID(ctx, user.CoupleID)
	if err != nil {
		return 0, err
	}
	return couple.ID, nil
}

// VerifyEmail 验证用户邮箱
func (s *userService) VerifyEmail(ctx context.Context, email, code string) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrEmailNotFound
		}
		return fmt.Errorf("查询用户时发生错误: %w", err)
	}

	// 验证验证码
	verified, err := s.emailSvc.VerifyCode(ctx, email, code)
	if err != nil {
		return fmt.Errorf("验证验证码时发生错误: %w", err)
	}

	if !verified {
		return ErrInvalidVerificationCode
	}

	// 验证成功后发送欢迎邮件
	err = s.emailSvc.SendWelcomeEmail(ctx, email, user.Username)
	if err != nil {
		// 记录错误但不影响验证流程
		fmt.Printf("发送欢迎邮件失败: %v", err)
	}

	return nil
}

// ForgotPassword 处理忘记密码请求
func (s *userService) ForgotPassword(ctx context.Context, email string) (string, error) {
	// 检查用户是否存在
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrEmailNotFound
		}
		return "", fmt.Errorf("查询用户时发生错误: %w", err)
	}

	// 生成重置令牌
	resetToken := generateResetToken()

	// 发送密码重置邮件
	err = s.emailSvc.SendPasswordResetEmail(ctx, email, resetToken)
	if err != nil {
		return "", fmt.Errorf("发送密码重置邮件失败: %w", err)
	}

	return resetToken, nil
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, email, token, newPassword string) error {
	// 验证重置令牌
	verified, err := s.emailSvc.VerifyPasswordResetToken(ctx, email, token)
	if err != nil {
		return fmt.Errorf("验证重置令牌时发生错误: %w", err)
	}

	if !verified {
		return ErrInvalidResetToken
	}

	// 获取用户
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrEmailNotFound
		}
		return fmt.Errorf("查询用户时发生错误: %w", err)
	}

	// 对新密码进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// GenerateVerificationCode 生成6位验证码
func (s *userService) GenerateVerificationCode() string {
	const codeLength = 6
	const charset = "0123456789"

	b := make([]byte, codeLength)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果随机数生成失败，使用时间戳作为备选
			n := time.Now().UnixNano() % int64(len(charset))
			b[i] = charset[n]
			continue
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}

// ResendVerificationCode 重发验证码
func (s *userService) ResendVerificationCode(ctx context.Context, email string) (string, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrEmailNotFound
		}
		return "", fmt.Errorf("查询用户时发生错误: %w", err)
	}

	// 生成新的验证码
	verificationCode := s.GenerateVerificationCode()

	// 发送验证码
	err = s.emailSvc.SendVerificationEmail(ctx, email, user.Username, verificationCode)
	if err != nil {
		return "", fmt.Errorf("发送验证码失败: %w", err)
	}

	return verificationCode, nil
}

// generateResetToken 生成密码重置令牌
func generateResetToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return time.Now().String()
	}
	return hex.EncodeToString(b)
}
