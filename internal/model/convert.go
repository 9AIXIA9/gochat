package model

import (
	"gochat/internal/domain"
)

func (u *GormUser) ToDomain() *domain.User {
	return &domain.User{
		Number:  u.Number,
		Name:    u.Name,
		PwdHash: u.PwdHash,
	}
}

func UserFromDomain(user *domain.User) *GormUser {
	return &GormUser{
		Number:  user.Number,
		Name:    user.Name,
		PwdHash: user.PwdHash,
	}
}

func (r *GormRoom) ToDomain() *domain.Room {
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

func RoomFromDomain(room *domain.Room) *GormRoom {
	return &GormRoom{
		Number:       room.Number,
		Name:         room.Name,
		Owner:        room.Owner,
		SecretHash:   room.SecretHash,
		Description:  room.Description,
		CurrentUsers: room.CurrentUsers,
		MaxUsers:     room.MaxUsers,
	}
}

func (m *GormMessage) ToDomain() domain.Message {
	return domain.Message{
		UserNumber: m.UserNumber,
		RoomNumber: m.RoomNumber,
		Content:    m.Content,
		SendTime:   m.SentAt,
	}
}

func MessageFromDomain(message domain.Message) *GormMessage {
	return &GormMessage{
		UserNumber: message.UserNumber,
		RoomNumber: message.RoomNumber,
		Content:    message.Content,
		SentAt:     message.SendTime,
	}
}
