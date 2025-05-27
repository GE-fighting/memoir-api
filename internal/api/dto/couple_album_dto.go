package dto

type CreateCoupleAlbumRequest struct {
	UserID      int64  `json:"user_id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
}

type UpdateCoupleAlbumRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
