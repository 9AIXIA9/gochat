package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"gochat/internal/config"
	"gochat/internal/domain"
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

func GenerateUserNumber() domain.UserNumber {
	return domain.UserNumber(node.Generate().Int64())
}

func GenerateRoomNumber() domain.RoomNumber {
	return domain.RoomNumber(node.Generate().Int64())
}
