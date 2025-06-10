package models

import (
	"encoding/json"
	"time"
)

// TimelineEvent 时间轴事件
type TimelineEvent struct {
	Base
	CoupleID    int64           `json:"couple_id,string" gorm:"not null"`
	EventDate   time.Time       `json:"event_date" gorm:"type:date;not null"`
	Title       string          `json:"title" gorm:"type:varchar(100);not null"`
	Description json.RawMessage `json:"description,omitempty" gorm:"type:jsonb"`
	LocationID  *int64          `json:"location_id,string,omitempty"`

	// 关联 - 没有外键约束
	Couple       Couple       `json:"-" gorm:"-"`
	Location     *Location    `json:"location,omitempty" gorm:"-"`
	PhotosVideos []PhotoVideo `json:"photos_videos,omitempty" gorm:"-"`
}
