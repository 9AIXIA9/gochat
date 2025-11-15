package snowflake

import (
	"gochat/internal/social/application"
	"gochat/internal/social/domain"

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

func (g *RoomNumberGenerator) Generate() domain.RoomNumber {
	return domain.RoomNumber(g.node.Generate().String())
}
