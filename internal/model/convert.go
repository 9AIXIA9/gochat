package model

import (
	"gochat/internal/domain"
)

func (u *User) ToDomain() *domain.User {
	return domain.NewUser(u.Number, u.Name, u.PwdHash)
}

func UserFromDomain(user *domain.User) *User {
	return &User{
		Number:  user.Number(),
		Name:    user.Name(),
		PwdHash: user.PwdHash(),
	}
}

func (r *Room) ToDomain() *domain.Room {
	return domain.NewRoom(r.Number, r.Name, r.SecretHash, r.Description, r.CurrentUsers, r.MaxUsers, r.Owner)
}

func RoomFromDomain(room *domain.Room) *Room {
	return &Room{
		Number:       room.Number(),
		Name:         room.Name(),
		Owner:        room.Owner(),
		SecretHash:   room.SecretHash(),
		Description:  room.Description(),
		CurrentUsers: room.CurrentUsers(),
		MaxUsers:     room.MaxUsers(),
	}
}

func (m *Message) ToDomain() *domain.Message {
	return domain.NewMessage(m.From, m.To, m.Content, m.SentAt)
}

func MessageFromDomain(message *domain.Message) *Message {
	return &Message{
		From:    message.From(),
		To:      message.To(),
		Content: message.Content(),
		SentAt:  message.SendAt(),
	}
}

func NewUserRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) *UserRoom {
	return &UserRoom{
		UserNumber: userNumber,
		RoomNumber: roomNumber,
	}
}
