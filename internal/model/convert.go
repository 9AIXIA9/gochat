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

func (m *Chat) ToDomain() *domain.Chat {
	return domain.NewChat(m.UserNumber, m.RoomNumber, m.Content, m.SentAt)
}

func ChatFromDomain(chat *domain.Chat) *Chat {
	return &Chat{
		UserNumber: chat.UserNumber(),
		RoomNumber: chat.RoomNumber(),
		Content:    chat.Content(),
		SentAt:     chat.SendTime(),
	}
}

func NewUserRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) *UserRoom {
	return &UserRoom{
		UserNumber: userNumber,
		RoomNumber: roomNumber,
	}
}
