package utils

import (
	"gochat/internal/domain"
	"gochat/internal/infra/repository"
)

// ModelToDomain 将数据库模型转换为领域模型
func ModelToDomain(from interface{}) interface{} {
	switch v := from.(type) {
	case repository.User:
		return domain.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case *repository.User:
		if v == nil {
			return nil
		}
		return domain.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case repository.Room:
		return domain.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case *repository.Room:
		if v == nil {
			return nil
		}
		return domain.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case repository.Chat:
		return domain.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case *repository.Chat:
		if v == nil {
			return nil
		}
		return domain.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case []repository.User:
		users := make([]domain.User, len(v))
		for i, u := range v {
			users[i] = ModelToDomain(u).(domain.User)
		}
		return users
	case []repository.Room:
		rooms := make([]domain.Room, len(v))
		for i, r := range v {
			rooms[i] = ModelToDomain(r).(domain.Room)
		}
		return rooms
	case []repository.Chat:
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
		return repository.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case *domain.User:
		if v == nil {
			return nil
		}
		return repository.User{
			Number:  v.Number,
			Name:    v.Name,
			PwdHash: v.PwdHash,
		}
	case domain.Room:
		return repository.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case *domain.Room:
		if v == nil {
			return nil
		}
		return repository.Room{
			Name:       v.Name,
			Number:     v.Number,
			SecretHash: v.SecretHash,
		}
	case domain.Chat:
		return repository.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case *domain.Chat:
		if v == nil {
			return nil
		}
		return repository.Chat{
			UserNumber: v.UserNumber,
			RoomNumber: v.RoomNumber,
			Content:    v.Content,
			SendTime:   v.SendTime,
		}
	case []domain.User:
		users := make([]repository.User, len(v))
		for i, u := range v {
			users[i] = DomainToModel(u).(repository.User)
		}
		return users
	case []domain.Room:
		rooms := make([]repository.Room, len(v))
		for i, r := range v {
			rooms[i] = DomainToModel(r).(repository.Room)
		}
		return rooms
	case []domain.Chat:
		chats := make([]repository.Chat, len(v))
		for i, c := range v {
			chats[i] = DomainToModel(c).(repository.Chat)
		}
		return chats
	default:
		return nil
	}
}
