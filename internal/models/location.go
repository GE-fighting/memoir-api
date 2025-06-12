package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

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

	// 关联表
	TimelineEventLocations []TimelineEventLocation `json:"-" gorm:"-"`
}
