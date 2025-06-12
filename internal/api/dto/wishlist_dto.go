package dto

import (
	"memoir-api/internal/models"
	"time"
)

// CreateWishlistRequest 创建心愿清单请求
type CreateWishlistRequest struct {
	CoupleID     int64   `json:"couple_id,string" binding:"required"`
	Title        string  `json:"title" binding:"required,max=100"`
	Description  string  `json:"description,omitempty"`
	Priority     int     `json:"priority" binding:"omitempty,min=1,max=3"` // 1-高，2-中，3-低
	Type         int     `json:"type" binding:"omitempty,min=1,max=2"`     // 1-日常，2-旅行
	ReminderDate *string `json:"reminder_date,omitempty"`                  // 格式: "2006-01-02"
}

// UpdateWishlistRequest 更新心愿清单请求
type UpdateWishlistRequest struct {
	ID           int64   `json:"id,string" binding:"required"`
	Title        string  `json:"title" binding:"omitempty,max=100"`
	Description  string  `json:"description,omitempty"`
	Priority     int     `json:"priority" binding:"omitempty,min=1,max=3"`
	Type         int     `json:"type" binding:"omitempty,min=1,max=2"` // 1-日常，2-旅行
	ReminderDate *string `json:"reminder_date,omitempty"`              // 格式: "2006-01-02"
}

// UpdateWishlistStatusRequest 更新心愿清单状态请求
type UpdateWishlistStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed"`
}

// WishlistDTO 心愿清单响应DTO
type WishlistDTO struct {
	ID           int64     `json:"id,string"`
	CoupleID     int64     `json:"couple_id,string"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Priority     int       `json:"priority"`
	Status       string    `json:"status"`
	Type         int       `json:"type"`                    // 1-日常，2-旅行
	ReminderDate *string   `json:"reminder_date,omitempty"` // 格式: "2006-01-02"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WishlistQueryParams 心愿清单查询参数
type WishlistQueryParams struct {
	PaginationRequest
	CoupleID int64  `form:"couple_id,string" binding:"required"`
	Status   string `form:"status" binding:"omitempty,oneof=pending completed"`
	Priority int    `form:"priority" binding:"omitempty,min=1,max=3"`
	Type     int    `form:"type" binding:"omitempty,min=1,max=2"` // 1-日常，2-旅行
}

// ToModel 将创建请求转换为模型对象
func (r *CreateWishlistRequest) ToModel() (*models.Wishlist, error) {
	wishlist := &models.Wishlist{
		CoupleID:    r.CoupleID,
		Title:       r.Title,
		Description: r.Description,
		Priority:    r.Priority,
		Type:        r.Type,
		Status:      "pending", // 默认状态为pending
	}

	// 如果Priority为0，设置默认值为2（中等优先级）
	if wishlist.Priority == 0 {
		wishlist.Priority = 2
	}

	// 如果Type为0，设置默认值为1（日常）
	if wishlist.Type == 0 {
		wishlist.Type = 1
	}

	// 处理提醒日期
	if r.ReminderDate != nil && *r.ReminderDate != "" {
		reminderDate, err := time.Parse("2006-01-02", *r.ReminderDate)
		if err != nil {
			return nil, err
		}
		wishlist.ReminderDate = &reminderDate
	}

	return wishlist, nil
}

// ApplyToModel 将更新请求应用到现有模型对象
func (r *UpdateWishlistRequest) ApplyToModel(wishlist *models.Wishlist) error {
	if r.Title != "" {
		wishlist.Title = r.Title
	}

	if r.Description != "" {
		wishlist.Description = r.Description
	}

	if r.Priority != 0 {
		wishlist.Priority = r.Priority
	}

	if r.Type != 0 {
		wishlist.Type = r.Type
	}

	// 处理提醒日期
	if r.ReminderDate != nil {
		if *r.ReminderDate == "" {
			// 空字符串表示清除提醒日期
			wishlist.ReminderDate = nil
		} else {
			reminderDate, err := time.Parse("2006-01-02", *r.ReminderDate)
			if err != nil {
				return err
			}
			wishlist.ReminderDate = &reminderDate
		}
	}

	return nil
}

// FromModel 从模型创建DTO
func WishlistFromModel(wishlist *models.Wishlist) WishlistDTO {
	dto := WishlistDTO{
		ID:          wishlist.ID,
		CoupleID:    wishlist.CoupleID,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		Priority:    wishlist.Priority,
		Status:      wishlist.Status,
		Type:        wishlist.Type,
		CreatedAt:   wishlist.CreatedAt,
		UpdatedAt:   wishlist.UpdatedAt,
	}

	// 处理提醒日期
	if wishlist.ReminderDate != nil {
		reminderDateStr := wishlist.ReminderDate.Format("2006-01-02")
		dto.ReminderDate = &reminderDateStr
	}

	return dto
}

// FromModels 从模型列表创建DTO列表
func WishlistsFromModels(wishlists []*models.Wishlist) []WishlistDTO {
	dtos := make([]WishlistDTO, 0, len(wishlists))
	for _, wishlist := range wishlists {
		dtos = append(dtos, WishlistFromModel(wishlist))
	}
	return dtos
}
