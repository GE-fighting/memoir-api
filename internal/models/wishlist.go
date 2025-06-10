package models

import (
	"time"
)

// Wishlist 心愿清单
type Wishlist struct {
	Base
	CoupleID     int64      `json:"couple_id,string" gorm:"not null"`
	Title        string     `json:"title" gorm:"type:varchar(100);not null"`
	Description  string     `json:"description,omitempty" gorm:"type:text"`
	Priority     int        `json:"priority" gorm:"not null;default:2"`                        // 1-高，2-中，3-低
	Status       string     `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // 'pending' or 'completed'
	ReminderDate *time.Time `json:"reminder_date,omitempty" gorm:"type:date"`

	// 关联 - 没有外键约束
	Couple Couple `json:"-" gorm:"-"`
}
