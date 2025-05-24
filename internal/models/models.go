package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

// Couple 情侣关系，包含设置字段
type Couple struct {
	Base
	AutoGenerateVideo     bool   `json:"auto_generate_video" gorm:"not null;default:true"`
	ReminderNotifications bool   `json:"reminder_notifications" gorm:"not null;default:true"`
	PairToken             string `json:"pair_token" gorm:"type:varchar(50);uniqueIndex;not null"`

	// 关联 - 没有外键约束
	Users []User `json:"users,omitempty" gorm:"-"`
}

// User 用户信息，包含深色模式字段
type User struct {
	Base
	CoupleID     int64  `json:"couple_id,string" gorm:""`
	Username     string `json:"username" gorm:"type:varchar(50);not null;index:users_email_key"`
	Email        string `json:"email" gorm:"type:varchar(100);not null;index:users_username_key"`
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null"`
	DarkMode     bool   `json:"dark_mode" gorm:"not null;default:false"`

	// 关联已移除
}

// BeforeUpdate GORM 更新用户前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 什么都不做，这样可以防止 GORM 尝试删除不存在的约束
	return nil
}

// Point 表示地理坐标点
type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// GormDataType 实现自定义类型的GORM接口
func (Point) GormDataType() string {
	return "geometry"
}

// GormDBDataType 实现GORM中的数据库特定类型
func (Point) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "geometry(point,4326)"
}

// Value 实现driver.Valuer接口，将结构转换为数据库值
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Lng, p.Lat), nil
}

// Scan 实现sql.Scanner接口，从数据库值读取到结构
func (p *Point) Scan(src interface{}) error {
	// 简化实现，实际项目中应使用专业的地理库如 github.com/twpayne/go-geom
	// 此处仅为示例代码
	if src == nil {
		p.Lng, p.Lat = 0, 0
		return nil
	}

	// 假设src是字符串，实际可能是二进制格式
	// 这里简化处理，实际需要根据PostgreSQL的EWKB格式正确解析
	var lng, lat float64
	_, err := fmt.Sscanf(src.(string), "POINT(%f %f)", &lng, &lat)
	if err != nil {
		return err
	}

	p.Lng, p.Lat = lng, lat
	return nil
}

// Location 地点
type Location struct {
	Base
	CoupleID    int64           `json:"couple_id,string" gorm:"not null"`
	Name        string          `json:"name" gorm:"type:varchar(100);not null"`
	Coordinates Point           `json:"coordinates" gorm:"type:geometry(Point,4326);not null"`
	Description json.RawMessage `json:"description,omitempty" gorm:"type:jsonb"`

	// 关联 - 没有外键约束
	Couple         Couple          `json:"-" gorm:"-"`
	TimelineEvents []TimelineEvent `json:"timeline_events,omitempty" gorm:"-"`
	PhotosVideos   []PhotoVideo    `json:"photos_videos,omitempty" gorm:"-"`
}

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

// PhotoVideo 照片和视频
type PhotoVideo struct {
	Base
	CoupleID     int64           `json:"couple_id,string" gorm:"not null"`
	MediaURL     string          `json:"media_url" gorm:"type:text;not null"`
	MediaType    string          `json:"media_type" gorm:"type:varchar(10);not null"` // 'photo' or 'video'
	Category     string          `json:"category" gorm:"type:varchar(50);not null"`   // 'date', 'travel', 'daily'
	ThumbnailURL *string         `json:"thumbnail_url,omitempty" gorm:"type:text"`
	Description  json.RawMessage `json:"description,omitempty" gorm:"type:jsonb"`
	EventID      *int64          `json:"event_id,string,omitempty"`
	LocationID   *int64          `json:"location_id,string,omitempty"`

	// 关联 - 没有外键约束
	Couple   Couple         `json:"-" gorm:"-"`
	Event    *TimelineEvent `json:"event,omitempty" gorm:"-"`
	Location *Location      `json:"location,omitempty" gorm:"-"`
}

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

// Wishlist 心愿清单
type Wishlist struct {
	Base
	CoupleID     int64           `json:"couple_id,string" gorm:"not null"`
	Title        string          `json:"title" gorm:"type:varchar(100);not null"`
	Description  json.RawMessage `json:"description,omitempty" gorm:"type:jsonb"`
	Priority     int             `json:"priority" gorm:"not null;default:2"`                        // 1-高，2-中，3-低
	Status       string          `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // 'pending' or 'completed'
	ReminderDate *time.Time      `json:"reminder_date,omitempty" gorm:"type:date"`

	// 关联 - 没有外键约束
	Couple Couple `json:"-" gorm:"-"`
}
