package uuid

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"

	"github.com/google/uuid"
)

var _ domain.RoomIDGenerator = (*RoomIDGenerator)(nil)

type RoomIDGenerator struct{}

func NewRoomIDGenerator() *RoomIDGenerator {
	return &RoomIDGenerator{}
}

func (RoomIDGenerator) Generate() kernel.RoomID {
	return kernel.RoomID(uuid.Must(uuid.NewV7()).String())
}
