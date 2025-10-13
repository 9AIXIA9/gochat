package model

import (
	"gochat/internal/domain"
)

func (u *User) ToDomain() *domain.User {
	return domain.NewUser(u.Number, u.PwdHash)
}

func UserFromDomain(user *domain.User) *User {
	return &User{
		Number:  user.Number(),
		PwdHash: user.PwdHash(),
	}
}

func (r *Room) ToDomain() *domain.Room {
	return domain.NewRoom(r.Number, r.SecretHash, r.CurrentUsers, r.MaxUsers, r.Owner)
}

func RoomFromDomain(room *domain.Room) *Room {
	return &Room{
		Number:       room.Number(),
		Owner:        room.Owner(),
		SecretHash:   room.SecretHash(),
		CurrentUsers: room.CurrentUsers(),
		MaxUsers:     room.MaxUsers(),
	}
}

func (m *Message) ToDomain() *domain.Message {
	return domain.NewMessage(m.ID, m.Sender, m.Recipient, m.Content, m.SentAt, m.Type)
}

func ToDomainMessages(msgs []*Message) []*domain.Message {
	domainMsgs := make([]*domain.Message, 0, len(msgs))

	for _, msg := range msgs {
		domainMsgs = append(domainMsgs, msg.ToDomain())
	}

	return domainMsgs
}

func MessageFromDomain(message *domain.Message) *Message {
	return &Message{
		ID:        message.ID(),
		Sender:    message.From(),
		Recipient: message.To(),
		Content:   message.Content(),
		SentAt:    message.SendAt(),
	}
}

func NewUserRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) *UserRoom {
	return &UserRoom{
		UserNumber: userNumber,
		RoomNumber: roomNumber,
	}
}
