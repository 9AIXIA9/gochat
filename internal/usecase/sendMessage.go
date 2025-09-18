package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/manager"
)

type SendMessage struct {
	repo    domain.MessageRepository
	manager *manager.Manager
}

func NewSendMessage(manager *manager.Manager, repo domain.MessageRepository) domain.SendMessageUsecase {
	return &SendMessage{
		repo:    repo,
		manager: manager,
	}
}

func (uc *SendMessage) Execute(ctx context.Context, req *domain.SendMessageRequest) (*domain.Response, error) {
	msg := domain.CreateMessage(req.AuthInfo.UserNumber, req.To, req.Content, req.SentAt, false)

	//尝试发送给客户端
	err := uc.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	//存储 message
	if err := uc.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	return domain.DefaultResponse, nil
}

func (uc *SendMessage) SaveMessage(ctx context.Context, msg *domain.Message) error {
	return uc.repo.Save(ctx, msg)
}

func (uc *SendMessage) Send(ctx context.Context, msg *domain.Message) error {
	return uc.manager.Send(ctx, msg)
}
