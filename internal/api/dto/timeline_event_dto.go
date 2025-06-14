package dto

import (
	"memoir-api/internal/models"
	"time"
)

// CreateTimelineEventRequest 创建时间线事件的请求
type CreateTimelineEventRequest struct {
	CoupleID      int64      `json:"couple_id,string" binding:"required"`
	StartDate     string     `json:"start_date" binding:"required"` // 格式：2006-01-02
	EndDate       string     `json:"end_date" binding:"required"`
	Title         string     `json:"title" binding:"required,max=100"`
	Content       string     `json:"content" binding:"required"`
	CoverURL      string     `json:"cover_url,omitempty"`
	LocationIDs   Int64Array `json:"location_ids,omitempty"`
	PhotoVideoIDs Int64Array `json:"photo_video_ids,omitempty"`
}

// UpdateTimelineEventRequest 更新时间线事件的请求
type UpdateTimelineEventRequest struct {
	Title         string     `json:"title,omitempty" binding:"omitempty,max=100"`
	StartDate     string     `json:"start_date" binding:"required"` // 格式：2006-01-02
	EndDate       string     `json:"end_date" binding:"required"`
	Content       string     `json:"content,omitempty"`
	LocationIDs   Int64Array `json:"location_ids,omitempty"`
	PhotoVideoIDs Int64Array `json:"photo_video_ids,omitempty"`
}

// TimelineEventQueryParams 查询时间线事件的参数
type TimelineEventQueryParams struct {
	PaginationRequest
	CoupleID   int64  `form:"couple_id,string" binding:"required"`
	StartDate  string `form:"start_date,omitempty"` // 格式：2006-01-02
	EndDate    string `form:"end_date,omitempty"`   // 格式：2006-01-02
	Title      string `form:"title,omitempty"`
	LocationID int64  `form:"location_id,string,omitempty"`
}

// ToModel 将创建请求转换为模型
func (r *CreateTimelineEventRequest) ToModel() (*models.TimelineEvent, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return nil, err
	}

	return &models.TimelineEvent{
		CoupleID:  r.CoupleID,
		StartDate: startDate,
		EndDate:   endDate,
		Title:     r.Title,
		Content:   r.Content,
		CoverURL:  r.CoverURL,
	}, nil
}

// ApplyToModel 将更新请求应用到模型
func (r *UpdateTimelineEventRequest) ApplyToModel(event *models.TimelineEvent) error {
	if r.StartDate != "" {
		eventDate, err := time.Parse("2006-01-02", r.StartDate)
		if err != nil {
			return err
		}
		event.StartDate = eventDate
	}
	if r.EndDate != "" {
		eventDate, err := time.Parse("2006-01-02", r.EndDate)
		if err != nil {
			return err
		}
		event.EndDate = eventDate
	}

	if r.Title != "" {
		event.Title = r.Title
	}

	if r.Content != "" {
		event.Content = r.Content
	}

	return nil
}

// TimelineEventResponse 时间线事件响应
type TimelineEventResponse struct {
	ID           int64               `json:"id,string"`
	CoupleID     int64               `json:"couple_id,string"`
	StartDate    string              `json:"start_date" binding:"required"` // 格式：2006-01-02
	EndDate      string              `json:"end_date" binding:"required"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	Locations    []models.Location   `json:"locations,omitempty"`
	PhotosVideos []models.PhotoVideo `json:"photos_videos,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// FromModel 从模型创建响应
func TimelineEventResponseFromModel(event *models.TimelineEvent) TimelineEventResponse {
	return TimelineEventResponse{
		ID:           event.ID,
		CoupleID:     event.CoupleID,
		StartDate:    event.StartDate.Format("2006-01-02"),
		EndDate:      event.EndDate.Format("2006-01-02"),
		Title:        event.Title,
		Content:      event.Content,
		Locations:    event.Locations,
		PhotosVideos: event.PhotosVideos,
		CreatedAt:    event.CreatedAt,
		UpdatedAt:    event.UpdatedAt,
	}
}
