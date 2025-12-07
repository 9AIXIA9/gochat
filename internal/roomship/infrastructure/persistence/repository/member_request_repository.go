package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/infrastructure/persistence/model"
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
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(memberRequest)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), memberRequest.GetEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (repo *MemberRequestRepository) Update(ctx context.Context, memberRequest *domain.MemberRequest) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(repo.toModel(memberRequest)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), memberRequest.GetEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (repo *MemberRequestRepository) FindByID(ctx context.Context, requestID kernel.OperationID) (*domain.MemberRequest, error) {
	var request model.MemberRequest
	if err := repo.db.WithContext(ctx).Where("id = ?", requestID).First(&request).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&request), nil
}

func (repo *MemberRequestRepository) ExistByUserIDAndRoomIDAndState(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID, state domain.MemberRequestState) (bool, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&model.MemberRequest{}).Where("applicant_id = ? AND room_id = ? AND state = ?", userID, roomID, state).Count(&count).Error; err != nil {
		return false, gormutils.TranslateError(err)
	}
	return count > 0, nil
}

func (repo *MemberRequestRepository) toModel(memberRequest *domain.MemberRequest) *model.MemberRequest {
	return &model.MemberRequest{
		ID:          memberRequest.ID(),
		ApplicantID: memberRequest.ApplicantID(),
		RoomID:      memberRequest.RoomID(),
		Content:     memberRequest.Content(),
		State:       memberRequest.State(),
		SentAt:      memberRequest.CreatedAt(),
		OperatorID:  memberRequest.OperatorID(),
		OperatedAt:  memberRequest.OperatedAt(),
	}
}

func (repo *MemberRequestRepository) toDomain(memberRequest *model.MemberRequest) *domain.MemberRequest {
	return domain.LoadMemberRequest(
		memberRequest.ID,
		memberRequest.State,
		memberRequest.ApplicantID,
		memberRequest.RoomID,
		memberRequest.Content,
		memberRequest.OperatorID,
		memberRequest.OperatedAt,
		memberRequest.SentAt,
	)
}
