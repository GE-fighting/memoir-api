package dto

import (
	"memoir-api/internal/models"
	"time"
)

// CreateAttachmentRequest 创建附件请求
type CreateAttachmentRequest struct {
	UserID    int64  `json:"user_id"`
	FileName  string `json:"file_name" binding:"required"`
	FileType  string `json:"file_type" binding:"required"`
	FileSize  int    `json:"file_size" binding:"required"`
	Url       string `json:"url" binding:"required"`
	CoupleID  int64  `json:"couple_id,string,omitempty"`
	SpaceType string `json:"space_type" binding:"required,oneof=personal couple"`
}

// AttachmentQueryParams 查询附件的参数
type AttachmentQueryParams struct {
	PaginationRequest
	UserID    int64  `form:"user_id"`
	CoupleID  int64  `form:"couple_id,string"`
	SpaceType string `form:"space_type" binding:"omitempty,oneof=personal couple"`
	FileType  string `form:"file_type"`
}

// AttachmentResponse 附件响应
type AttachmentResponse struct {
	ID        int64     `json:"id,string"`
	UserID    int64     `json:"user_id,string"`
	CoupleID  int64     `json:"couple_id,string,omitempty"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"file_type"`
	FileSize  int       `json:"file_size"`
	Url       string    `json:"url"`
	SpaceType string    `json:"space_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToModel 将创建请求转换为模型
func (r *CreateAttachmentRequest) ToModel() *models.Attachment {
	return &models.Attachment{
		UserID:    r.UserID,
		FileName:  r.FileName,
		FileType:  r.FileType,
		FileSize:  r.FileSize,
		Url:       r.Url,
		CoupleID:  r.CoupleID,
		SpaceType: r.SpaceType,
	}
}

// AttachmentFromModel 从模型创建响应DTO
func AttachmentFromModel(attachment *models.Attachment) AttachmentResponse {
	return AttachmentResponse{
		ID:        attachment.ID,
		UserID:    attachment.UserID,
		CoupleID:  attachment.CoupleID,
		FileName:  attachment.FileName,
		FileType:  attachment.FileType,
		FileSize:  attachment.FileSize,
		Url:       attachment.Url,
		SpaceType: attachment.SpaceType,
		CreatedAt: attachment.CreatedAt,
		UpdatedAt: attachment.UpdatedAt,
	}
}

// AttachmentsFromModels 从模型列表创建响应DTO列表
func AttachmentsFromModels(attachments []models.Attachment) []AttachmentResponse {
	result := make([]AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		result[i] = AttachmentFromModel(&attachment)
	}
	return result
}
