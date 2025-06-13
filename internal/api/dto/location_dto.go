package dto

import (
	"memoir-api/internal/models"
	"time"
)

// CreateLocationRequest 创建地点请求
type CreateLocationRequest struct {
	CoupleID    int64   `json:"couple_id,string" binding:"required"`
	Name        string  `json:"name" binding:"required,max=100"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Description string  `json:"description,omitempty"`
}

// UpdateLocationRequest 更新地点请求
type UpdateLocationRequest struct {
	Name        string  `json:"name,omitempty" binding:"omitempty,max=100"`
	Longitude   float64 `json:"longitude,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Description string  `json:"description,omitempty"`
}

// LocationQueryParams 查询地点的参数
type LocationQueryParams struct {
	PaginationRequest
	CoupleID int64  `form:"couple_id,string" binding:"required"`
	Name     string `form:"name,omitempty"`
}

// LocationResponse 地点响应
type LocationResponse struct {
	ID          int64     `json:"id,string"`
	CoupleID    int64     `json:"couple_id,string"`
	Name        string    `json:"name"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToModel 将创建请求转换为模型
func (r *CreateLocationRequest) ToModel() *models.Location {
	return &models.Location{
		CoupleID:    r.CoupleID,
		Name:        r.Name,
		Longitude:   r.Longitude,
		Latitude:    r.Latitude,
		Description: r.Description,
	}
}

// FromModel 从模型创建响应DTO
func LocationFromModel(location *models.Location) LocationResponse {
	return LocationResponse{
		ID:          location.ID,
		CoupleID:    location.CoupleID,
		Name:        location.Name,
		Longitude:   location.Longitude,
		Latitude:    location.Latitude,
		Description: location.Description,
		CreatedAt:   location.CreatedAt,
		UpdatedAt:   location.UpdatedAt,
	}
}

// FromModels 从模型列表创建响应DTO列表
func LocationsFromModels(locations []*models.Location) []LocationResponse {
	responses := make([]LocationResponse, len(locations))
	for i, location := range locations {
		responses[i] = LocationFromModel(location)
	}
	return responses
}
