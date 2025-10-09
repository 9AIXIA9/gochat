package usecase

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/types"
)

type userConnected struct {
	repo    domain.MessageRepository
	manager *manager.Manager
}

func NewUserConnected(repo domain.MessageRepository, m *manager.Manager) domain.UserConnectedUsecase {
	return &userConnected{
		repo:    repo,
		manager: m,
	}
}

func (uc *userConnected) Execute(ctx context.Context, number domain.UserNumber) error {
	msgs, err := uc.QueryUnsentMessages(ctx, number)
	if err != nil {
		return err
	}

	if len(msgs) == 0 {
		return nil
	}

	msgIDs := uc.SendUserManyMsgs(number, msgs)

	err = uc.UpdateMessagesSentToOneUser(ctx, number, msgIDs)
	if err != nil {
		return err
	}
	return nil
}

func (uc *userConnected) QueryUnsentMessages(ctx context.Context, number domain.UserNumber) ([]*domain.Message, error) {
	return uc.repo.QueryUnsentMessages(ctx, number)
}

func (uc *userConnected) SendUserManyMsgs(number domain.UserNumber, msgs []*domain.Message) []domain.MessageID {
	msgIDs := make([]domain.MessageID, 0, len(msgs))
	for _, msg := range msgs {
		err := uc.manager.SendMessage(number, msg)
		if err == nil {
			msgIDs = append(msgIDs, msg.ID())
		}
		if err != nil && !errors.Is(err, types.ErrNotFound) {
			zap.L().Error("send message failed",
				zap.Int64("user_number", int64(number)),
				zap.String("message_id", string(msg.ID())),
				zap.Error(err))
		}
	}
	return msgIDs
}

func (uc *userConnected) UpdateMessagesSentToOneUser(ctx context.Context, number domain.UserNumber, msgIDs []domain.MessageID) error {
	return uc.repo.UpdateMessagesSentToOneUser(ctx, number, msgIDs)
}
