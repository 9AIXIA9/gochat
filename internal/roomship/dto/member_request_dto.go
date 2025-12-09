package dto

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type MemberRequest struct {
	ID          kernel.OperationID        `json:"id"`
	State       domain.MemberRequestState `json:"state"`
	ApplicantID kernel.UserID             `json:"applicant_id"`
	RoomID      kernel.RoomID             `json:"room_id"`
	Content     string                    `json:"content"`
	OperatorID  kernel.UserID             `json:"operator_id"`
	OperatedAt  time.Time                 `json:"operated_at"`
	CreatedAt   time.Time                 `json:"created_at"`
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
