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
	msg := domain.CreateMessage(req.UserNumber, req.To, req.Content, req.SentAt)

	//存储 message并查询应该发送的用户
	userNumbers, err := uc.SaveAndQueryUserNumberShouldSent(ctx, msg)
	if err != nil {
		return nil, err
	}

	//发送给用户
	numbersSent := uc.SendMsgToManyUsers(msg, userNumbers)

	//存储已发送的信息
	if err := uc.UpdateMessageSentToManyUsers(ctx, msg.ID(), numbersSent); err != nil {
		return nil, err
	}

	return domain.DefaultResponse, nil
}

func (uc *SendMessage) SendMsgToManyUsers(msg *domain.Message, numbers []domain.UserNumber) []domain.UserNumber {
	return uc.manager.SendMsgToManyUsers(msg, numbers)
}

func (uc *SendMessage) SaveAndQueryUserNumberShouldSent(ctx context.Context, msg *domain.Message) ([]domain.UserNumber, error) {
	return uc.repo.SaveAndQueryUserNumberShouldSent(ctx, msg)
}

func (uc *SendMessage) UpdateMessageSentToManyUsers(ctx context.Context, msgID domain.MessageID, userNumbers []domain.UserNumber) error {
	return uc.repo.UpdateMessageSentToManyUsers(ctx, msgID, userNumbers)
}
