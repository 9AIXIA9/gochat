package dto

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Roomship struct {
	ID        domain.RoomshipID   `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	RoomID    kernel.RoomID       `json:"room_id" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	UserID    kernel.UserID       `json:"user_id" example:"019b593b-462e-74d6-bfda-0e103a172193"`
	Role      domain.RoomshipRole `json:"role" example:"member"`
	CreatedAt time.Time           `json:"created_at" example:"2025-12-26 05:38:19.740"`
}

func ToRoomshipDTO(roomship *domain.Roomship) *Roomship {
	return &Roomship{
		ID:        roomship.ID(),
		RoomID:    roomship.RoomID(),
		UserID:    roomship.UserID(),
		Role:      roomship.Role(),
		CreatedAt: roomship.CreatedAt(),
	}
}

func ToRoomshipDTOs(roomships []*domain.Roomship) []*Roomship {
	dtos := make([]*Roomship, 0, len(roomships))
	for _, roomship := range roomships {
		dtos = append(dtos, ToRoomshipDTO(roomship))
	}
	return dtos
}
