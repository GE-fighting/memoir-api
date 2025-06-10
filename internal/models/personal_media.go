package models

// PersonalMedia 个人空间的照片、视频和其他内容
type PersonalMedia struct {
	Base
	UserID       int64   `json:"user_id,string" gorm:"not null;index"` // 关键区别：属于单个用户
	MediaURL     string  `json:"media_url" gorm:"type:text;not null"`
	MediaType    string  `json:"media_type" gorm:"type:varchar(10);not null"` // 'photo' or 'video'
	Category     *string `json:"category" gorm:"type:varchar(50)"`            // 'photos', 'videos', 'notes', 'favorites'
	ThumbnailURL *string `json:"thumbnail_url,omitempty" gorm:"type:text"`
	Description  *string `json:"description,omitempty" gorm:"type:text"`
	Title        *string `json:"title" gorm:"type:varchar(100)"`

	// 关联
	User User `json:"-" gorm:"-"`
}
