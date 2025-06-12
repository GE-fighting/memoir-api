package models

// TimelineEventPhotoVideo 时间线事件与照片/视频的关联表
type TimelineEventPhotoVideo struct {
	Base
	TimelineEventID int64 `json:"timeline_event_id,string" gorm:"not null;index"`
	PhotoVideoID    int64 `json:"photo_video_id,string" gorm:"not null;index"`

	// 关联 - 没有外键约束，但在应用层面维护关系
	TimelineEvent TimelineEvent `json:"-" gorm:"-"`
	PhotoVideo    PhotoVideo    `json:"-" gorm:"-"`
}
