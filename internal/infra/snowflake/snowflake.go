package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"gochat/internal/config"
	"gochat/internal/domain"
)

var _ domain.NumberGenerator = (*NumberGenerator)(nil)

type NumberGenerator struct {
	node *snowflake.Node
}

func NewNumberGenerator(conf *config.Snowflake) (*NumberGenerator, error) {
	node, err := snowflake.NewNode(conf.Node)
	if err != nil {
		return nil, err
	}

	return &NumberGenerator{node: node}, nil
}

func (g *NumberGenerator) GenerateNumber() domain.BaseNumber {
	return domain.BaseNumber(g.node.Generate().Int64())
}
