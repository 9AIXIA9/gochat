package repository

import (
	"context"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.MemberRequestRepository = (*MemberRequestRepository)(nil)

type MemberRequestRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewMemberRequestRepository(db *gorm.DB, eventRepo event.Repository) *MemberRequestRepository {
	return &MemberRequestRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *MemberRequestRepository) Create(ctx context.Context, memberRequest *domain.MemberRequest) error {
	//TODO implement me
	panic("implement me")
}

func (repo *MemberRequestRepository) Update(ctx context.Context, memberRequest *domain.MemberRequest) error {
	//TODO implement me
	panic("implement me")
}

func (repo *MemberRequestRepository) FindByID(ctx context.Context, requestID kernel.OperationID) (*domain.MemberRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *MemberRequestRepository) ExistByUserIDAndRoomIDAndState(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID, state domain.MemberRequestState) (bool, error) {
	//TODO implement me
	panic("implement me")
}
