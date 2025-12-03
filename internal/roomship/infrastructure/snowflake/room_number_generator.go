package snowflake

import (
	"gochat/internal/roomship/domain"

	"github.com/bwmarrin/snowflake"
)

var _ domain.RoomNumberGenerator = (*RoomNumberGenerator)(nil)

type RoomNumberGenerator struct {
	node *snowflake.Node
}

func NewRoomNumberGenerator(machineNode int64) (*RoomNumberGenerator, error) {
	node, err := snowflake.NewNode(machineNode)
	if err != nil {
		return nil, err
	}

	return &RoomNumberGenerator{node: node}, nil
}

func (g *RoomNumberGenerator) Generate() domain.RoomNumber {
	return domain.RoomNumber(g.node.Generate().String())
}
