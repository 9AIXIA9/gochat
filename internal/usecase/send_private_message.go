package usecase

import (
	"context"
	"errors"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/types"
)

type SendPrivateMessage struct {
	repo    domain.MessageRepository
	manager *manager.Manager
}

func NewSendPrivateMessage(manager *manager.Manager, repo domain.MessageRepository) domain.SendPrivateMessageUsecase {
	return &SendPrivateMessage{
		repo:    repo,
		manager: manager,
	}
}

func (uc *SendPrivateMessage) Execute(ctx context.Context, req *domain.SendMessageRequest) (*domain.Response, error) {
	msg := domain.CreateMessage(req.UserNumber, req.To, req.Content, req.SentAt, domain.MessageTypePrivate)

	//存储 message
	if err := uc.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	//发送给用户
	err := uc.SendMessage(msg)
	if err == nil {
		// 存储已发送的信息
		if err := uc.UpdateMessageSent(ctx, req.UserNumber, msg.ID()); err != nil {
			return nil, err
		}
		return domain.NewSuccessResponse(domain.SendMessageResponse{
			MessageID: msg.ID(),
		}), nil
	}

	if errors.Is(err, types.ErrNotFound) {
		return domain.NewSuccessResponse(domain.SendMessageResponse{
			MessageID: msg.ID(),
		}), nil
	}

	return nil, err
}

func (uc *SendPrivateMessage) SaveMessage(ctx context.Context, msg *domain.Message) error {
	return uc.repo.Save(ctx, msg)
}

func (uc *SendPrivateMessage) SendMessage(msg *domain.Message) error {
	return uc.manager.SendMessage(domain.UserNumber(msg.To()), msg)
}

func (uc *SendPrivateMessage) UpdateMessageSent(ctx context.Context, userNumber domain.UserNumber, msgID domain.MessageID) error {
	return uc.repo.UpdateMessageSent(ctx, userNumber, msgID)
}
