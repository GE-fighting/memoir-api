package dto

import "memoir-api/internal/models"

type DashboardDTO struct {
	StoryCount int                `json:"story_count"`
	MediaCount int                `json:"media_count"`
	AlbumCount int                `json:"album_count"`
	CoupleDays int                `json:"couple_days"`
	Locations  []*models.Location `json:"locations,omitempty"`
}
