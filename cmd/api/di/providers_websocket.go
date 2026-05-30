package di

import (
	chatDomain "gochat/internal/chat/domain"
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
	provideAllowedAction,
	provideGatewayService,
	provideGatewayIngressHandler,
)

type AllowedAction map[event.Topic]struct{}

func provideGatewayManager() *core.Manager {
	return core.NewManager(1024)
}

func provideGatewayService(manager *core.Manager) contract.GatewayService {
	return local.NewLocalGatewayService(manager)
}

func provideAllowedAction() AllowedAction {
	return map[event.Topic]struct{}{
		chatDomain.TopicSendPrivateMessageCommand: {},
		chatDomain.TopicSendRoomMessageCommand:    {},
	}
}

func provideGatewayUpstreamHandler(publisher event.SyncPublisher, generator event.IDGenerator, allowedAction AllowedAction) contract.UpstreamHandler {
	return gateway.NewUpstreamRouter(publisher, generator, allowedAction)
}

func provideGatewayIngressHandler(hub *core.Manager, upstream contract.UpstreamHandler, sessionIDGenerator *uuid.SessionIDGenerator, checkOrigin checkOrigin) *gatewayWebsocket.IngressHandler {
	return gatewayWebsocket.NewIngressHandler(hub, upstream, checkOrigin, sessionIDGenerator)
}
