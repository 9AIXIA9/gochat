package snowflake

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"

	"github.com/bwmarrin/snowflake"
)

var _ domain.UserNumberGenerator = (*UserNumberGenerator)(nil)

type UserNumberGenerator struct {
	node *snowflake.Node
}

func NewUserNumberGenerator(machineNode int64) (*UserNumberGenerator, error) {
	node, err := snowflake.NewNode(machineNode)
	if err != nil {
		return nil, err
	}

	return &UserNumberGenerator{node: node}, nil
}

func (g *UserNumberGenerator) Generate() kernel.UserNumber {
	return kernel.UserNumber(g.node.Generate().String())
}
