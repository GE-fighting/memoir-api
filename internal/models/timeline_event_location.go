package models

// TimelineEventLocation 时间线事件与位置的关联表
type TimelineEventLocation struct {
	Base
	TimelineEventID int64 `json:"timeline_event_id,string" gorm:"not null;index"`
	LocationID      int64 `json:"location_id,string" gorm:"not null;index"`

	// 关联 - 没有外键约束，但在应用层面维护关系
	TimelineEvent TimelineEvent `json:"-" gorm:"-"`
	Location      Location      `json:"-" gorm:"-"`
}
