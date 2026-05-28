package di

import (
	chatApp "gochat/internal/chat/application"
	"gochat/internal/delivery/ws_gateway"
	validatorInfra "gochat/internal/infrastructure/validator" // 导入实际验证器类型
	gatewayWebsocket "gochat/internal/ws_gateway/adapter/websocket"
	"gochat/internal/ws_gateway/api/local"
	"gochat/internal/ws_gateway/core"

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
) *ws_gateway.UpstreamRouter {
	return ws_gateway.NewUpstreamRouter(validator, chatSendPrivateMessage, chatSendRoomMessage).(*ws_gateway.UpstreamRouter)
}

func provideGatewayIngressHandler(hub *core.Manager, upstream *ws_gateway.UpstreamRouter) *gatewayWebsocket.IngressHandler {
	return gatewayWebsocket.NewIngressHandler(hub, upstream)
}
