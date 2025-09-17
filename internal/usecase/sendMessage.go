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
	msg := domain.NewMessage(req.AuthInfo.UserNumber, req.To, req.Content, req.SentAt)

	if err := uc.SaveMessage(ctx, msg); err != nil {

	}

	ok, err := uc.Send(ctx, msg)
	if err != nil {

	}

	if ok {
		//成功发送 创建消息发送事件
		return
	}
	return
}

func (uc *SendMessage) SaveMessage(ctx context.Context, msg *domain.Message) error {
	return uc.repo.Save(ctx, msg)
}

func (uc *SendMessage) Send(ctx context.Context, msg *domain.Message) (bool, error) {
	return uc.manager.Send(ctx, msg)
}
