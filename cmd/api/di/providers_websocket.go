package di

import (
	"gochat/internal/delivery/gateway"
	gatewayWebsocket "gochat/internal/gateway/adapter/websocket"
	"gochat/internal/gateway/api/local"
	"gochat/internal/gateway/core"
	"gochat/internal/gateway/infrastructure/uuid"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"

	"github.com/google/wire"
)

var WebsocketSet = wire.NewSet(
	provideGatewayManager,
	provideGatewayUpstreamHandler,
	provideGatewayService,
	provideGatewayIngressHandler,
)

func provideGatewayManager() *core.Manager {
	return core.NewManager(1024)
}

func provideGatewayService(manager *core.Manager) contract.GatewayService {
	return local.NewLocalGatewayService(manager)
}

func provideGatewayUpstreamHandler(publisher event.SyncPublisher, generator event.IDGenerator) contract.UpstreamHandler {
	// 定义外网 Action 到内网 Kafka Topic 的静态映射表
	allowedActions := map[string]event.Topic{
		// 客户端 action -> 内部领域事件 Topic
		"send_private_message": event.Topic("chat.send_private_message"),
		"send_room_message":    event.Topic("chat.send_room_message"),
		// 后续新增业务直接在这里白名单注册
	}
	return gateway.NewUpstreamRouter(publisher, generator, allowedActions)
}

func provideGatewayIngressHandler(hub *core.Manager, upstream contract.UpstreamHandler, sessionIDGenerator *uuid.SessionIDGenerator, checkOrigin checkOrigin) *gatewayWebsocket.IngressHandler {
	return gatewayWebsocket.NewIngressHandler(hub, upstream, checkOrigin, sessionIDGenerator)
}
