package snowflake

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"

	"github.com/bwmarrin/snowflake"
)

var _ application.UserNumberGenerator = (*NumberGenerator)(nil)

type NumberGenerator struct {
	node *snowflake.Node
}

func NewNumberGenerator(machineNode int64) (*NumberGenerator, error) {
	node, err := snowflake.NewNode(machineNode)
	if err != nil {
		return nil, err
	}

	return &NumberGenerator{node: node}, nil
}

func (g *NumberGenerator) Generate() domain.UserNumber {
	return domain.UserNumber(g.node.Generate().String())
}
