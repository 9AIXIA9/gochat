package usecase

import (
	"context"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/manager"
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

	uc.SendMessages(ctx, msgs)

	err = uc.UpdateMessagesSent(ctx, msgs)
	if err != nil {
		return err
	}
	return nil
}

func (uc *userConnected) QueryUnsentMessages(ctx context.Context, number domain.UserNumber) ([]*domain.Message, error) {
	return uc.repo.QueryUnsentMessages(ctx, number)
}

func (uc *userConnected) SendMessages(ctx context.Context, msgs []*domain.Message) {
	for _, msg := range msgs {
		if err := uc.manager.Send(ctx, msg); err != nil {
			zap.L().Error("send message failed", zap.Error(err))
		}
	}
}

func (uc *userConnected) UpdateMessagesSent(ctx context.Context, msgs []*domain.Message) error {
	return uc.repo.UpdateMessagesSent(ctx, msgs)
}
