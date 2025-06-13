package models

// Location 地点
type Location struct {
	Base
	CoupleID    int64  `json:"couple_id,string" gorm:"not null"`
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Longitude   float64 `json:"longitude" gorm:"not null;index"`
	Latitude    float64 `json:"latitude" gorm:"not null;index"`
	Description string `json:"description,omitempty" gorm:"type:text"`

	// 关联 - 没有外键约束
	Couple         Couple          `json:"-" gorm:"-"`
	TimelineEvents []TimelineEvent `json:"timeline_events,omitempty" gorm:"-"`
	PhotosVideos   []PhotoVideo    `json:"photos_videos,omitempty" gorm:"-"`

	// 关联表
	TimelineEventLocations []TimelineEventLocation `json:"-" gorm:"-"`
}
