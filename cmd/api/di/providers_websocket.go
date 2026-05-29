package di

import (
	chatApp "gochat/internal/chat/application"
	"gochat/internal/delivery/gateway"
	gatewayWebsocket "gochat/internal/gateway/adapter/websocket"
	"gochat/internal/gateway/api/local"
	"gochat/internal/gateway/core"
	"gochat/internal/gateway/infrastructure/uuid"
	validatorInfra "gochat/internal/infrastructure/validator"
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

func provideGatewayUpstreamHandler(
	validator *validatorInfra.Validator,
	chatSendPrivateMessage chatApp.SendPrivateMessageUseCase,
	chatSendRoomMessage chatApp.SendRoomMessageUseCase,
) *gateway.UpstreamRouter {
	return gateway.NewUpstreamRouter(validator, chatSendPrivateMessage, chatSendRoomMessage).(*gateway.UpstreamRouter)
}

func provideGatewayIngressHandler(hub *core.Manager, upstream *gateway.UpstreamRouter, sessionIDGenerator *uuid.SessionIDGenerator, checkOrigin checkOrigin) *gatewayWebsocket.IngressHandler {
	return gatewayWebsocket.NewIngressHandler(hub, upstream, checkOrigin, sessionIDGenerator)
}
