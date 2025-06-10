package models

// CoupleAlbum 情侣相册
type CoupleAlbum struct {
	Base
	CoupleID     int64        `json:"couple_id,string" gorm:"not null"`
	Title        string       `json:"title" gorm:"type:varchar(100);not null"`
	Description  string       `json:"description,omitempty" gorm:"type:text"`
	CoverURL     *string      `json:"cover_url,omitempty" gorm:"type:text"`
	Count        int          `json:"count" gorm:"not null;default:0"`
	PhotosVideos []PhotoVideo `json:"photos_videos,omitempty" gorm:"-"`
}
