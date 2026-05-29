package di

import (
	gatewayWebsocket "gochat/internal/gateway/adapter/websocket"
	"gochat/internal/gateway/api/local"
	"gochat/internal/gateway/core"
	"gochat/internal/gateway/infrastructure/uuid"
	"gochat/internal/shared/contract"

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

func provideGatewayUpstreamHandler() contract.UpstreamHandler {
	return nil
}

func provideGatewayIngressHandler(hub *core.Manager, upstream contract.UpstreamHandler, sessionIDGenerator *uuid.SessionIDGenerator, checkOrigin checkOrigin) *gatewayWebsocket.IngressHandler {
	return gatewayWebsocket.NewIngressHandler(hub, upstream, checkOrigin, sessionIDGenerator)
}
