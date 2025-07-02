package model

import (
	"gochat/internal/domain"
)

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		Number:  u.Number,
		Name:    u.Name,
		PwdHash: u.PwdHash,
	}
}

func UserFromDomain(user *domain.User) *User {
	return &User{
		Number:  user.Number,
		Name:    user.Name,
		PwdHash: user.PwdHash,
	}
}

func (r *Room) ToDomain() *domain.Room {
	return &domain.Room{
		Name:         r.Name,
		Number:       r.Number,
		SecretHash:   r.SecretHash,
		Description:  r.Description,
		CurrentUsers: r.CurrentUsers,
		MaxUsers:     r.MaxUsers,
		Owner:        r.Owner,
	}
}

func RoomFromDomain(room *domain.Room) *Room {
	return &Room{
		Number:       room.Number,
		Name:         room.Name,
		Owner:        room.Owner,
		SecretHash:   room.SecretHash,
		Description:  room.Description,
		CurrentUsers: room.CurrentUsers,
		MaxUsers:     room.MaxUsers,
	}
}

func (m *Chat) ToDomain() domain.Chat {
	return domain.Chat{
		UserNumber: m.UserNumber,
		RoomNumber: m.RoomNumber,
		Content:    m.Content,
		SendTime:   m.SentAt,
	}
}

func ChatFromDomain(chat domain.Chat) *Chat {
	return &Chat{
		UserNumber: chat.UserNumber,
		RoomNumber: chat.RoomNumber,
		Content:    chat.Content,
		SentAt:     chat.SendTime,
	}
}

func NewUserRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) *UserRoom {
	return &UserRoom{
		UserNumber: userNumber,
		RoomNumber: roomNumber,
	}
}
