package uuid

import (
	"gochat/internal/roomship/domain"

	"github.com/google/uuid"
)

var _ domain.RoomshipIDGenerator = (*RoomshipIDGenerator)(nil)

type RoomshipIDGenerator struct{}

func NewRoomshipIDGenerator() *RoomshipIDGenerator {
	return &RoomshipIDGenerator{}
}

func (RoomshipIDGenerator) Generate() domain.RoomshipID {
	return domain.RoomshipID(uuid.Must(uuid.NewV7()).String())
}
