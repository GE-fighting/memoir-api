package models

import (
	"time"
)

// TimelineEvent 时间轴事件
type TimelineEvent struct {
	Base
	CoupleID  int64     `json:"couple_id,string" gorm:"not null"`
	EventDate time.Time `json:"event_date" gorm:"type:date;not null"`
	Title     string    `json:"title" gorm:"type:varchar(100);not null"`
	Content   string    `json:"content,omitempty" gorm:"type:text;not null"`
	CoverURL  string    `json:"cover_url,omitempty" gorm:"type:text"`
	// 关联 - 没有外键约束
	Couple       Couple       `json:"-" gorm:"-"`
	Locations    []Location   `json:"locations,omitempty" gorm:"-"`
	PhotosVideos []PhotoVideo `json:"photos_videos,omitempty" gorm:"-"`

	// 关联表
	TimelineEventLocations    []TimelineEventLocation   `json:"-" gorm:"-"`
	TimelineEventPhotosVideos []TimelineEventPhotoVideo `json:"-" gorm:"-"`
}
