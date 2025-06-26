package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"gochat/internal/config"
	"gochat/internal/domain"
)

var (
	node *snowflake.Node
)

func Init(config config.Snowflake) (err error) {
	if node, err = snowflake.NewNode(config.Node); err != nil {
		return err
	}
	return nil
}

func GenerateUserNumber() domain.UserNumber {
	return domain.UserNumber(node.Generate().Int64())
}

func GenerateRoomNumber() domain.RoomNumber {
	return domain.RoomNumber(node.Generate().Int64())
}
