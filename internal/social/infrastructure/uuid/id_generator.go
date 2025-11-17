package uuid

import (
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"

	"github.com/google/uuid"
)

var _ application.RoomIDGenerator = (*RoomIDGenerator)(nil)

type RoomIDGenerator struct{}

func NewRoomIDGenerator() *RoomIDGenerator {
	return &RoomIDGenerator{}
}

func (RoomIDGenerator) Generate() kernel.RoomID {
	return kernel.RoomID(uuid.Must(uuid.NewV7()).String())
}
