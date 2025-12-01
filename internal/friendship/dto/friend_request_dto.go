package dto

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type FriendRequest struct {
	ID      kernel.OperationID `json:"id"`
	From    kernel.UserID      `json:"from"`
	To      kernel.UserID      `json:"to"`
	Content string             `json:"content"`
	State   string             `json:"state"`
	SentAt  time.Time          `json:"sent_at"` //UTC
}

func ToFriendRequestDTO(request *domain.FriendRequest) *FriendRequest {
	return &FriendRequest{
		ID:      request.ID(),
		From:    request.From(),
		To:      request.To(),
		Content: request.Content(),
		State:   request.State().String(),
		SentAt:  request.SentAt(),
	}
}

func ToFriendRequestDTOs(requests []*domain.FriendRequest) []*FriendRequest {
	dtos := make([]*FriendRequest, 0, len(requests))
	for _, request := range requests {
		dtos = append(dtos, ToFriendRequestDTO(request))
	}
	return dtos
}
