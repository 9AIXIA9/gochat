package converter

import (
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"
)

var _ gormutils.GenericModelConverter[*model.Room, *domain.Room] = (*RoomConverter)(nil)

type RoomConverter struct {
}

func (c *RoomConverter) ToModel(room *domain.Room) *model.Room {
	return &model.Room{
		ID:                room.ID(),
		Number:            room.Number(),
		PasswordEncrypted: room.PasswordEncrypted(),
	}
}

func (c *RoomConverter) ToDomain(room *model.Room) *domain.Room {
	return domain.NewRoom(room.ID, room.Number, room.PasswordEncrypted)
}
