package dto

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type FriendRequest struct {
	ID      kernel.OperationID        `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	From    kernel.UserID             `json:"from" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	To      kernel.UserID             `json:"to" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	Content string                    `json:"content" example:"Hi, let's be friends!"`
	State   domain.FriendRequestState `json:"state" example:"pending"`
	SentAt  time.Time                 `json:"sent_at" example:"2025-12-26 05:38:19.740"` //UTC
}

func ToFriendRequestDTO(request *domain.FriendRequest) *FriendRequest {
	return &FriendRequest{
		ID:      request.ID(),
		From:    request.From(),
		To:      request.To(),
		Content: request.Content(),
		State:   request.State(),
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
