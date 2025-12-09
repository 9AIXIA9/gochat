package dto

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Friendship struct {
	ID        domain.FriendshipID `json:"id"`
	UserID1   kernel.UserID       `json:"user_id_1"`
	UserID2   kernel.UserID       `json:"user_id_2"`
	CreatedAt time.Time           `json:"created_at"`
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
