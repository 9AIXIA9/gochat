package di

import (
	chatApp "gochat/internal/chat/application"
	"gochat/internal/delivery/gateway"
	gatewayWebsocket "gochat/internal/gateway/adapter/websocket"
	"gochat/internal/gateway/api/local"
	"gochat/internal/gateway/core"
	validatorInfra "gochat/internal/infrastructure/validator" // 导入实际验证器类型

	"github.com/google/wire"
)

var WebsocketSet = wire.NewSet(
	provideGatewayManager,
	provideGatewayUpstreamHandler,
	provideGatewayService,
	provideGatewayIngressHandler,
)

func provideGatewayManager() *core.Manager {
	return core.NewManager()
}

func provideGatewayService(manager *core.Manager) *local.WSGatewayService {
	return local.NewWSGatewayService(manager).(*local.WSGatewayService)
}

func provideGatewayUpstreamHandler(
	validator *validatorInfra.Validator, // 改为实际结构
	chatSendPrivateMessage chatApp.SendPrivateMessageUseCase,
	chatSendRoomMessage chatApp.SendRoomMessageUseCase,
) *gateway.UpstreamRouter {
	return gateway.NewUpstreamRouter(validator, chatSendPrivateMessage, chatSendRoomMessage).(*gateway.UpstreamRouter)
}

func provideGatewayIngressHandler(hub *core.Manager, upstream *gateway.UpstreamRouter) *gatewayWebsocket.IngressHandler {
	return gatewayWebsocket.NewIngressHandler(hub, upstream)
}
