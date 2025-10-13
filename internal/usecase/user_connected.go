package usecase

import (
	"context"
	"gochat/internal/domain"
)

type userConnected struct {
	domain.MessageSender
	domain.SendUnsentMessageAggregate
}

func NewUserConnected(sender domain.MessageSender,
	sendUnsentMessageAggregate domain.SendUnsentMessageAggregate) domain.UserConnectedUsecase {
	return &userConnected{
		MessageSender:              sender,
		SendUnsentMessageAggregate: sendUnsentMessageAggregate,
	}
}

func (uc *userConnected) Execute(ctx context.Context, number domain.UserNumber) error {
	msgs, err := uc.FindUnsentMessages(ctx, number)
	if err != nil {
		return err
	}

	if len(msgs) == 0 {
		return nil
	}

	for _, msg := range msgs {
		//忽略发送出错的消息
		if err := uc.SendMessage(number, msg); err == nil {
			if err := uc.UpdateMessageSent(ctx, number, msg.ID()); err != nil {

			}
		}
	}

	return nil
}
