package models

// PhotoVideo 照片和视频
type PhotoVideo struct {
	Base
	CoupleID     int64  `json:"couple_id,string" gorm:"not null"`
	AlbumID      int64  `json:"album_id,string" gorm:"not null"`
	MediaURL     string `json:"media_url" gorm:"type:text;not null"`
	MediaType    string `json:"media_type" gorm:"type:varchar(10);not null"` // 'photo' or 'video'
	ThumbnailURL string `json:"thumbnail_url,omitempty" gorm:"type:text"`
	Description  string `json:"description,omitempty" gorm:"type:text"`
	Title        string `json:"title,omitempty" gorm:"type:varchar(100)"`

	// 关联 - 没有外键约束
	Couple         Couple          `json:"-" gorm:"-"`
	TimelineEvents []TimelineEvent `json:"timeline_events,omitempty" gorm:"-"`
	Location       *Location       `json:"location,omitempty" gorm:"-"`

	// 关联表
	TimelineEventPhotosVideos []TimelineEventPhotoVideo `json:"-" gorm:"-"`
}
