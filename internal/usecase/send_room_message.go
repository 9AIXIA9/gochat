package usecase

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/types"
)

type SendRoomMessage struct {
	repo    domain.MessageRepository
	manager *manager.Manager
}

func NewSendRoomMessage(manager *manager.Manager, repo domain.MessageRepository) domain.SendRoomMessageUsecase {
	return &SendRoomMessage{
		repo:    repo,
		manager: manager,
	}
}

func (uc *SendRoomMessage) Execute(ctx context.Context, req *domain.SendMessageRequest) (*domain.Response, error) {
	msg := domain.CreateMessage(req.UserNumber, req.To, req.Content, req.SentAt, domain.MessageTypeRoom)

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

	return domain.NewSuccessResponse(domain.SendMessageResponse{
		MessageID: msg.ID(),
	}), nil
}

func (uc *SendRoomMessage) SendMsgToManyUsers(msg *domain.Message, numbers []domain.UserNumber) []domain.UserNumber {
	numbersSent := make([]domain.UserNumber, 0, len(numbers))
	for _, number := range numbers {
		err := uc.manager.SendMessage(number, msg)
		if err == nil {
			numbersSent = append(numbersSent, number)
		}
		if err != nil && !errors.Is(err, types.ErrNotFound) {
			zap.L().Error("send message failed",
				zap.Int64("user_number", int64(number)),
				zap.String("message_id", string(msg.ID())),
				zap.Error(err))
		}
	}
	return numbersSent
}

func (uc *SendRoomMessage) SaveAndQueryUserNumberShouldSent(ctx context.Context, msg *domain.Message) ([]domain.UserNumber, error) {
	return uc.repo.SaveAndQueryUserNumberShouldSent(ctx, msg)
}

func (uc *SendRoomMessage) UpdateMessageSentToManyUsers(ctx context.Context, msgID domain.MessageID, userNumbers []domain.UserNumber) error {
	return uc.repo.UpdateMessageSentToManyUsers(ctx, msgID, userNumbers)
}
