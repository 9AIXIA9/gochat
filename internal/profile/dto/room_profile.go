package dto

import (
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomProfile struct {
	ID           kernel.RoomID `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	Name         string        `json:"name" example:"we're family"`
	Introduction string        `json:"introduction" example:"This is our family room."`
	CreatedAt    time.Time     `json:"created_at" example:"2025-12-26 05:38:19.740"`
}

func ToRoomProfileDTO(profile *domain.RoomProfile) *RoomProfile {
	return &RoomProfile{
		ID:           profile.ID(),
		Name:         profile.Name(),
		Introduction: profile.Introduction(),
		CreatedAt:    profile.CreatedAt(),
	}
}
