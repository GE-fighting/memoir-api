package models

type Attachment struct {
	Base
	FileName  string `json:"file_name" gorm:"type:varchar(255);not null"`
	FileType  string `json:"file_type" gorm:"type:varchar(10);not null"`
	FileSize  int    `json:"file_size" gorm:"not null"`
	Url       string `json:"url" gorm:"type:varchar(255);not null"`
	UserID    int64  `json:"user_id" gorm:"not null"`
	CoupleID  int64  `json:"couple_id"`
	SpaceType string `json:"space_type" gorm:"type:varchar(20);not null;default:'personal'"`
}
