package dto

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Roomship struct {
	ID        domain.RoomshipID   `json:"id"`
	RoomID    kernel.RoomID       `json:"room_id"`
	UserID    kernel.UserID       `json:"user_id"`
	Role      domain.RoomshipRole `json:"role"`
	CreatedAt time.Time           `json:"created_at"`
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
