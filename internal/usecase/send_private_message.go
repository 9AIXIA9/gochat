package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/utils"
)

type SendPrivateMessage struct {
	domain.MessageSender
	domain.SendMessageAggregate
}

func NewSendPrivateMessage(sender domain.MessageSender, sendMessageAggregate domain.SendMessageAggregate) domain.SendPrivateMessageUsecase {
	return &SendPrivateMessage{
		SendMessageAggregate: sendMessageAggregate,
		MessageSender:        sender,
	}
}

func (uc *SendPrivateMessage) Execute(ctx context.Context, req *domain.SendMessageRequest) (*domain.Response, error) {
	msg := domain.CreateMessage(req.UserNumber, req.To, req.Content, req.SentAt, domain.MessageTypePrivate)

	//存储 message
	if err := uc.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	//发送给用户
	err := uc.SendMessage(domain.UserNumber(msg.To()), msg)
	if err == nil {
		// 存储已发送的信息
		if err := uc.UpdateMessageSent(ctx, req.UserNumber, msg.ID()); err != nil {
			return nil, err
		}
		return domain.NewSuccessResponse(domain.SendMessageResponse{
			MessageID: msg.ID(),
		}), nil
	}

	if utils.IsNotFound(err) {
		return domain.NewSuccessResponse(domain.SendMessageResponse{
			MessageID: msg.ID(),
		}), nil
	}

	return nil, err
}
