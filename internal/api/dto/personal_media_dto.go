package dto

// CreatePersonalMediaWithURLRequest 通过URL创建个人媒体请求
type CreatePersonalMediaWithURLRequest struct {
	UserID       int64  `json:"userID"`
	MediaType    string `json:"media_type" binding:"required,oneof=photo video"`
	Category     string `json:"category"`
	Title        string `json:"title"`
	MediaURL     string `json:"media_url" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url" binding:"required"`
	Description  string `json:"description"`
}

// QueryPersonalMediaRequest 查询个人媒体请求
type QueryPersonalMediaRequest struct {
	PaginationRequest
	UserID    int64  `json:"user_id"`    // 用户ID
	Category  string `json:"category"`   // 分类
	MediaType string `json:"media_type"` // 媒体类型
}
