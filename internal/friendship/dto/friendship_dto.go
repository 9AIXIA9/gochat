package dto

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Friendship struct {
	ID        domain.FriendshipID `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	UserID1   kernel.UserID       `json:"user_id_1" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	UserID2   kernel.UserID       `json:"user_id_2" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	CreatedAt time.Time           `json:"created_at" example:"2025-12-26 05:38:19.740"`
}

func ToFriendshipDTO(friendship *domain.Friendship) *Friendship {
	return &Friendship{
		ID:        friendship.ID(),
		UserID1:   friendship.UserID1(),
		UserID2:   friendship.UserID2(),
		CreatedAt: friendship.CreatedAt(),
	}
}

func ToFriendshipDTOs(friendships []*domain.Friendship) []*Friendship {
	dtos := make([]*Friendship, 0, len(friendships))
	for _, friendship := range friendships {
		dtos = append(dtos, ToFriendshipDTO(friendship))
	}
	return dtos
}
