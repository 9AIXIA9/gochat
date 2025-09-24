package snowflake

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/utils/timeout"
	"log"
)

var (
	node *snowflake.Node
)

func MustInit(config *config.Snowflake) {
	var err error
	if node, err = snowflake.NewNode(config.Node); err != nil {
		log.Fatalf("init snowflake node failed : %v", err)
	}
}

func GenerateUserNumber(ctx context.Context) (domain.UserNumber, error) {
	return timeout.ConvertAndExecuteWithResponse(ctx, func() (domain.UserNumber, error) {
		return domain.UserNumber(node.Generate().Int64()), nil
	})
}

func GenerateRoomNumber(ctx context.Context) (domain.RoomNumber, error) {
	return timeout.ConvertAndExecuteWithResponse(ctx, func() (domain.RoomNumber, error) {
		return domain.RoomNumber(node.Generate().Int64()), nil
	})
}
