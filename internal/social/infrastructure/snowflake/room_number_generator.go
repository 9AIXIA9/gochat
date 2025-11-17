package snowflake

import (
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"

	"github.com/bwmarrin/snowflake"
)

var _ application.RoomNumberGenerator = (*RoomNumberGenerator)(nil)

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

func (g *RoomNumberGenerator) Generate() kernel.RoomNumber {
	return kernel.RoomNumber(g.node.Generate().String())
}
