package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// 初始化随机数生成器
var secureRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	ErrUserExists      = errors.New("用户已存在")
	ErrInvalidPassword = errors.New("密码不正确")
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
}

// userService 用户服务实现
type userService struct {
	*BaseService
	userRepo   repository.UserRepository
	coupleRepo repository.CoupleRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, coupleRepo repository.CoupleRepository) UserService {
	return &userService{
		BaseService: NewBaseService(userRepo),
		userRepo:    userRepo,
		coupleRepo:  coupleRepo,
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
		b[i] = charset[secureRand.Intn(len(charset))]
	}
	return string(b)
}
