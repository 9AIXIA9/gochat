package usecase

import (
	"context"
	"gochat/internal/domain"
)

type SendRoomMessage struct {
	domain.MessageSender
	domain.SendMessageAggregate
	domain.RoomMemberFinder
}

func NewSendRoomMessage(sender domain.MessageSender,
	sendMessageAggregate domain.SendMessageAggregate,
	finder domain.RoomMemberFinder) domain.SendRoomMessageUsecase {
	return &SendRoomMessage{
		MessageSender:        sender,
		SendMessageAggregate: sendMessageAggregate,
		RoomMemberFinder:     finder,
	}
}

func (uc *SendRoomMessage) Execute(ctx context.Context, req *domain.SendMessageRequest) (*domain.Response, error) {
	msg := domain.CreateMessage(req.UserNumber, req.To, req.Content, req.SentAt, domain.MessageTypeRoom)

	//存储 message
	if err := uc.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	//查询该房间用户
	userNumbers, err := uc.FindRoomMembers(ctx, domain.RoomNumber(msg.To()))
	if err != nil {
		return nil, err
	}

	//发送给每个用户
	for _, number := range userNumbers {
		//忽略发送出错
		if err := uc.SendMessage(number, msg); err == nil {
			if err := uc.UpdateMessageSent(ctx, number, msg.ID()); err != nil {
			}
		}
	}

	return domain.NewSuccessResponse(domain.SendMessageResponse{
		MessageID: msg.ID(),
	}), nil
}
