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

type CoupleInfoDTO struct {
	CoupleId        int64  `json:"couple_id,string"`
	CoupleName      string `json:"couple_name"`
	CoupleDays      int    `json:"couple_days"`
	AnniversaryDate string `json:"anniversary_date"`
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
