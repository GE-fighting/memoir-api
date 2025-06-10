package models

import (
	"time"

	"gorm.io/gorm"
)

// Base 包含所有模型共享的基础字段
type Base struct {
	ID        int64          `json:"id,string" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate 在创建记录前自动生成ID
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		b.ID = GenerateID()
	}
	return nil
}
