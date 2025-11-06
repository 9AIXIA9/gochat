package snowflake

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"

	"github.com/bwmarrin/snowflake"
)

var _ application.UserNumberGenerator = (*UserNumberGenerator)(nil)

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

func (g *UserNumberGenerator) Generate() domain.UserNumber {
	return domain.UserNumber(g.node.Generate().String())
}
