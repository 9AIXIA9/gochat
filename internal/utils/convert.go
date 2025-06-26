package utils

import (
	"gochat/internal/domain"
	"gochat/internal/infra/model"
)

// ModelToDomain 将数据库模型转换为领域模型
func ModelToDomain(from interface{}) interface{} {
	switch v := from.(type) {
	case model.User:
		return domain.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case *model.User:
		if v == nil {
			return nil
		}
		return domain.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case model.Room:
		return domain.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case *model.Room:
		if v == nil {
			return nil
		}
		return domain.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case model.Chat:
		return domain.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case *model.Chat:
		if v == nil {
			return nil
		}
		return domain.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case []model.User:
		users := make([]domain.User, len(v))
		for i, u := range v {
			users[i] = ModelToDomain(u).(domain.User)
		}
		return users
	case []model.Room:
		rooms := make([]domain.Room, len(v))
		for i, r := range v {
			rooms[i] = ModelToDomain(r).(domain.Room)
		}
		return rooms
	case []model.Chat:
		chats := make([]domain.Chat, len(v))
		for i, c := range v {
			chats[i] = ModelToDomain(c).(domain.Chat)
		}
		return chats
	default:
		return nil
	}
}

// DomainToModel 将领域模型转换为数据库模型
func DomainToModel(from interface{}) interface{} {
	switch v := from.(type) {
	case domain.User:
		return model.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case *domain.User:
		if v == nil {
			return nil
		}
		return model.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case domain.Room:
		return model.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case *domain.Room:
		if v == nil {
			return nil
		}
		return model.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case domain.Chat:
		return model.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case *domain.Chat:
		if v == nil {
			return nil
		}
		return model.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case []domain.User:
		users := make([]model.User, len(v))
		for i, u := range v {
			users[i] = DomainToModel(u).(model.User)
		}
		return users
	case []domain.Room:
		rooms := make([]model.Room, len(v))
		for i, r := range v {
			rooms[i] = DomainToModel(r).(model.Room)
		}
		return rooms
	case []domain.Chat:
		chats := make([]model.Chat, len(v))
		for i, c := range v {
			chats[i] = DomainToModel(c).(model.Chat)
		}
		return chats
	default:
		return nil
	}
}
