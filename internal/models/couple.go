package models

import (
	"time"
)

// Couple 情侣关系，包含设置字段
type Couple struct {
	Base
	AutoGenerateVideo     bool      `json:"auto_generate_video" gorm:"not null;default:true"`
	ReminderNotifications bool      `json:"reminder_notifications" gorm:"not null;default:true"`
	PairToken             string    `json:"pair_token" gorm:"type:varchar(50);uniqueIndex;not null"`
	AnniversaryDate       time.Time `json:"anniversary_date" gorm:"type:date"`
	// 关联 - 没有外键约束
	Users []User `json:"users,omitempty" gorm:"-"`
}
