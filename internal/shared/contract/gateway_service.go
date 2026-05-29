package contract

import (
	"context"
)

// GatewayService 是对外部业务系统暴露的推送接口。
type GatewayService interface {
	// PushToUser 向指定用户的所有在线设备推送数据
	PushToUser(ctx context.Context, userID string, payload []byte) error
}
