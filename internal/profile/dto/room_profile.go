package dto

import (
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomProfile struct {
	ID           kernel.RoomID `json:"id"`
	Name         string        `json:"name"`
	Introduction string        `json:"introduction"`
	CreatedAt    time.Time     `json:"created_at"`
}

func ToRoomProfileDTO(profile *domain.RoomProfile) *RoomProfile {
	return &RoomProfile{
		ID:           profile.ID(),
		Name:         profile.Name(),
		Introduction: profile.Introduction(),
		CreatedAt:    profile.CreatedAt(),
	}
}
