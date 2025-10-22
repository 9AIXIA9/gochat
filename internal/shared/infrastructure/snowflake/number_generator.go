package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"gochat/internal/shared/kernel"
)

var _ kernel.NumberGenerator = (*NumberGenerator)(nil)

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

func (g *NumberGenerator) Generate() kernel.Number {
	return kernel.Number(g.node.Generate().String())
}
