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
	chatReadPrivateMessages chatApp.ReadPrivateMessagesUseCase,
	chatReadRoomMessages chatApp.ReadRoomMessagesUseCase,
) *websocket.Router {
	router := websocket.NewRouter(validator)

	router.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewTelemetryMiddleware(appConfig.Name),
	)

	router.NoRoute(websocketDelivery.NewNotFoundHandler())

	{
		router.Handle(chatWebsocket.SendPrivateMessageTopic, chatWebsocket.NewSendPrivateMessageHandler(chatSendPrivateMessage))
		router.Handle(chatWebsocket.SendRoomMessageTopic, chatWebsocket.NewSendRoomMessageHandler(chatSendRoomMessage))
		router.Handle(chatWebsocket.ReadPrivateMessagesTopic, chatWebsocket.NewReadPrivateMessagesHandler(chatReadPrivateMessages))
		router.Handle(chatWebsocket.ReadRoomMessagesTopic, chatWebsocket.NewReadRoomMessagesHandler(chatReadRoomMessages))
	}
	return router
}
