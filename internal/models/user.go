package models

import (
	"gorm.io/gorm"
)

// User 用户信息，包含深色模式字段
type User struct {
	Base
	CoupleID     int64  `json:"couple_id,string" gorm:""`
	Username     string `json:"username" gorm:"type:varchar(50);not null;index:users_email_key"`
	Email        string `json:"email" gorm:"type:varchar(100);not null;index:users_username_key"`
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null"`
	DarkMode     bool   `json:"dark_mode" gorm:"not null;default:false"`

	// 关联已移除
}

// BeforeUpdate GORM 更新用户前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 什么都不做，这样可以防止 GORM 尝试删除不存在的约束
	return nil
}
