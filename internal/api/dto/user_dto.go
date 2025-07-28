package dto

import (
	"memoir-api/internal/models"
	"time"
)

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	DarkMode bool   `json:"dark_mode" binding:"omitempty"`
}

type UpdateUserPasswordDTO struct {
	UserID          int64  `json:"user_id"`
	CurrentPassword string `json:"current_password" binding:"omitempty"`
	NewPassword     string `json:"new_password" binding:"omitempty"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID        int64     `json:"id,string"`
	CoupleID  int64     `json:"couple_id,string"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	DarkMode  bool      `json:"dark_mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProfileResponse 用户个人资料响应（更详细的信息）
type UserProfileResponse struct {
	ID        int64     `json:"id,string"`
	CoupleID  int64     `json:"couple_id,string"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	DarkMode  bool      `json:"dark_mode"`
	HasCouple bool      `json:"has_couple"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel 从用户模型创建响应DTO
func UserFromModel(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		CoupleID:  user.CoupleID,
		Username:  user.Username,
		Email:     user.Email,
		DarkMode:  user.DarkMode,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserProfileFromModel 从用户模型创建个人资料响应DTO
func UserProfileFromModel(user *models.User) UserProfileResponse {
	return UserProfileResponse{
		ID:        user.ID,
		CoupleID:  user.CoupleID,
		Username:  user.Username,
		Email:     user.Email,
		DarkMode:  user.DarkMode,
		HasCouple: user.CoupleID > 0,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ApplyUpdates 将更新请求应用到用户模型
func (r *UpdateUserRequest) ApplyUpdates(user *models.User) {
	if r.Username != "" {
		user.Username = r.Username
	}
	if r.Email != "" {
		user.Email = r.Email
	}
	if r.DarkMode {
		user.DarkMode = true
	}
}
