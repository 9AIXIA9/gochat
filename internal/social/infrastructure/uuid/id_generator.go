package uuid

import (
	"gochat/internal/social/application"
	"gochat/internal/social/domain"

	"github.com/google/uuid"
)

var _ application.RoomIDGenerator = (*RoomIDGenerator)(nil)

type RoomIDGenerator struct{}

func NewRoomIDGenerator() *RoomIDGenerator {
	return &RoomIDGenerator{}
}

func (RoomIDGenerator) Generate() domain.RoomID {
	return domain.RoomID(uuid.Must(uuid.NewV7()).String())
}
