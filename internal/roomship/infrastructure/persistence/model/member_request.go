package model

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type MemberRequest struct {
	ID          kernel.OperationID        `gorm:"primaryKey;type:char(36)"`
	ApplicantID kernel.UserID             `gorm:"not null;index:idx_member_req_applicant_room_state,priority:1;index:idx_member_req_applicant;"`
	RoomID      kernel.RoomID             `gorm:"not null;index:idx_member_req_applicant_room_state,priority:2;index:idx_member_req_room_id;"`
	Content     string                    `gorm:"type:varchar(500)"`
	State       domain.MemberRequestState `gorm:"not null;index:idx_member_req_applicant_room_state,priority:3;index:idx_member_req_state"`
	SentAt      time.Time                 `gorm:"not null;index:idx_member_req_to_sentat,priority:2;index:idx_member_req_sentat"`
	OperatorID  kernel.UserID             `gorm:"not null;index"`
	OperatedAt  time.Time
}

func (*MemberRequest) TableName() string {
	return "roomship_member_requests"
}
