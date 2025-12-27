package dto

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type MemberRequest struct {
	ID          kernel.OperationID        `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	State       domain.MemberRequestState `json:"state" example:"pending"`
	ApplicantID kernel.UserID             `json:"applicant_id" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	RoomID      kernel.RoomID             `json:"room_id" example:"019b593b-462e-74d6-bfda-0e103a172193"`
	Content     string                    `json:"content" example:"I would like to join the room."`
	OperatorID  kernel.UserID             `json:"operator_id" example:"019b593b-462e-74d6-bfda-0e103a172194"`
	OperatedAt  time.Time                 `json:"operated_at" example:"2025-12-26 05:38:19.740"`
	CreatedAt   time.Time                 `json:"created_at" example:"2025-12-26 05:38:20.740"`
}

func ToMemberRequestDTO(request *domain.MemberRequest) *MemberRequest {
	return &MemberRequest{
		ID:          request.ID(),
		State:       request.State(),
		ApplicantID: request.ApplicantID(),
		RoomID:      request.RoomID(),
		Content:     request.Content(),
		OperatorID:  request.OperatorID(),
		OperatedAt:  request.OperatedAt(),
		CreatedAt:   request.CreatedAt(),
	}
}

func ToMemberRequestDTOs(requests []*domain.MemberRequest) []*MemberRequest {
	dtos := make([]*MemberRequest, 0, len(requests))
	for _, request := range requests {
		dtos = append(dtos, ToMemberRequestDTO(request))
	}
	return dtos
}
