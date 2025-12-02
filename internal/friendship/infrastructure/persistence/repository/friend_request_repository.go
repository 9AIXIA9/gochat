package repository

import (
	"context"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.FriendRequestRepository = (*FriendRequestRepository)(nil)

const (
	defaultLimit = 10
	maxLimit     = 100
)

type FriendRequestRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewFriendRequestRepository(db *gorm.DB, eventRepo event.Repository) *FriendRequestRepository {
	return &FriendRequestRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *FriendRequestRepository) Create(ctx context.Context, friendRequest *domain.FriendRequest) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(friendRequest)).Error; err != nil {
			return err
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), friendRequest.GetEvents()); err != nil {
			return err
		}

		return nil
	})
}

func (repo *FriendRequestRepository) Update(ctx context.Context, friendRequest *domain.FriendRequest) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(repo.toModel(friendRequest)).Error; err != nil {
			return err
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), friendRequest.GetEvents()); err != nil {
			return err
		}

		return nil
	})
}

func (repo *FriendRequestRepository) ExistByUserIDAndState(ctx context.Context, userID1, userID2 kernel.UserID, state domain.FriendRequestState) (bool, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&model.FriendRequest{}).
		Where("(`from` = ? AND `to` = ? AND state = ?) OR (`from` = ? AND `to` = ? AND state = ?)",
			userID1, userID2, state, userID2, userID1, state).
		Count(&count).Error; err != nil {
		return false, gormutils.TranslateError(err)
	}
	return count > 0, nil
}

func (repo *FriendRequestRepository) FindByID(ctx context.Context, requestID kernel.OperationID) (*domain.FriendRequest, error) {
	var request model.FriendRequest
	if err := repo.db.WithContext(ctx).First(&request, "id = ?", requestID).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	return repo.toDomain(&request), nil
}

func (repo *FriendRequestRepository) FindsByUserID(
	ctx context.Context,
	userID kernel.UserID,
	limit int,
	baseID ...kernel.OperationID,
) ([]*domain.FriendRequest, error) {
	switch {
	case limit <= 0:
		limit = defaultLimit
	case limit > maxLimit:
		limit = maxLimit
	}

	var models []model.FriendRequest

	// 基于发送时间的稳定倒序 + id 兜底，规避随机 UUID 无序的问题
	q := repo.db.WithContext(ctx).
		Model(&model.FriendRequest{}).
		Where(&model.FriendRequest{To: userID}). // 让 GORM 处理列名转义，避免 `to` 保留字
		Order("sent_at DESC").
		Order("id DESC").
		Limit(limit)

	if len(baseID) > 0 && baseID[0] != "" {
		var cursor model.FriendRequest
		if err := repo.db.WithContext(ctx).
			Select("id", "sent_at").
			Where("id = ? AND `to` = ?", baseID[0], userID).
			Take(&cursor).Error; err != nil {
			return nil, err
		}

		// 复合游标：(sent_at, id) < (cursor.sent_at, cursor.id)
		q = q.Where("(sent_at < ?) OR (sent_at = ? AND id < ?)", cursor.SentAt, cursor.SentAt, cursor.ID)
	}

	if err := q.Find(&models).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(models), nil
}

func (repo *FriendRequestRepository) toModel(request *domain.FriendRequest) *model.FriendRequest {
	return &model.FriendRequest{
		ID:      request.ID(),
		From:    request.From(),
		To:      request.To(),
		Content: request.Content(),
		State:   request.State(),
		SentAt:  request.SentAt(),
	}
}

func (repo *FriendRequestRepository) toDomain(request *model.FriendRequest) *domain.FriendRequest {
	return domain.LoadFriendRequest(
		request.ID,
		request.From,
		request.To,
		request.Content,
		request.State,
		request.SentAt,
	)
}

func (repo *FriendRequestRepository) toDomains(requests []model.FriendRequest) []*domain.FriendRequest {
	domainRequests := make([]*domain.FriendRequest, 0, len(requests))
	for _, r := range requests {
		domainRequests = append(domainRequests, repo.toDomain(&r))
	}
	return domainRequests
}
