package dto

import (
	"memoir-api/internal/models"
	"time"
)

type CreateCoupleRequest struct {
	UserID          int64  `json:"user_id"`
	PairToken       string `json:"pair_token" binding:"required"`
	AnniversaryDate string `json:"anniversary_date" binding:"required"`
}

func (r *CreateCoupleRequest) ToCouple() (models.Couple, error) {
	anniversaryDate, err := time.Parse("2006-01-02", r.AnniversaryDate)
	if err != nil {
		return models.Couple{}, err
	}

	return models.Couple{
		PairToken:       r.PairToken,
		AnniversaryDate: anniversaryDate,
	}, nil
}
