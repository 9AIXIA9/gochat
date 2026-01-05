package di

import (
	"gochat/config"
	chatApp "gochat/internal/chat/application"
	chatWebsocket "gochat/internal/chat/port/websocket"
	websocketDelivery "gochat/internal/delivery/websocket/handler"
	"gochat/internal/delivery/websocket/middleware"
	"gochat/internal/infrastructure/websocket"

	"github.com/google/wire"
	gorillaWebsocket "github.com/gorilla/websocket"
)

var WebsocketSet = wire.NewSet(
	provideWebsocketUpgrader,
	provideWebsocketManager,
	provideWebsocketRouter,
)

func provideWebsocketUpgrader(appConfig *config.App) *gorillaWebsocket.Upgrader {
	return websocket.NewUpgrader(appConfig.CORS.AllowOrigins)
}

func provideWebsocketManager() *websocket.Manager {
	return websocket.NewManager()
}

func provideWebsocketRouter(
	appConfig *config.App,
	validator websocket.Validator,
	chatSendPrivateMessage chatApp.SendPrivateMessageUseCase,
	chatSendRoomMessage chatApp.SendRoomMessageUseCase,

) *websocket.Router {
	router := websocket.NewRouter(validator)

	router.Use(
		middleware.NewRecoverMiddleware(),
		middleware.NewTraceMiddleware(appConfig.Name+".websocket"),
		middleware.NewLoggerMiddleware(),
	)

	if appConfig.Breaker != nil {
		appConfig.Breaker.Name = appConfig.Name + "_websocket_circuit_breaker"
		router.Use(middleware.NewCircuitBreakMiddleware(appConfig.Breaker))
	}

	router.NoRoute(websocketDelivery.NewNotFoundHandler())

	{
		router.Handle(chatWebsocket.SendPrivateMessageTopic, chatWebsocket.NewSendPrivateMessageHandler(chatSendPrivateMessage, validator))
		router.Handle(chatWebsocket.SendRoomMessageTopic, chatWebsocket.NewSendRoomMessageHandler(chatSendRoomMessage, validator))
	}
	return router
}
