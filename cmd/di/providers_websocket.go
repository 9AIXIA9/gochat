package di

import (
	"gochat/config"
	websocketDelivery "gochat/internal/delivery/websocket/handler"
	"gochat/internal/delivery/websocket/middleware"
	"gochat/internal/infrastructure/websocket"
	notificationApp "gochat/internal/notification/application"
	notificationWebsocket "gochat/internal/notification/port/websocket"

	"github.com/google/wire"
	gorillaWebsocket "github.com/gorilla/websocket"
)

var WebsocketSet = wire.NewSet(
	provideWebsocketUpgrader,
	provideWebsocketManager,
	provideWebsocketRouter,
	provideWebsocketServer,
)

func provideWebsocketUpgrader(appConfig *config.App) *gorillaWebsocket.Upgrader {
	return websocket.NewUpgrader(appConfig.CORS.AllowOrigins)
}

func provideWebsocketManager() *websocket.Manager {
	return websocket.NewManager()
}

func provideWebsocketRouter(
	appConfig *config.App,
	notificationReadPrivateMessage notificationApp.ReadPrivateMessageUseCase,
	notificationReadRoomMessage notificationApp.ReadRoomMessageUseCase,
) *websocket.Router {
	router := websocket.NewRouter()

	router.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewTelemetryMiddleware(appConfig.Name),
	)

	router.NoRoute(websocketDelivery.NewNotFoundHandler())

	router.Handle(notificationWebsocket.ReadPrivateMessageTopic, notificationWebsocket.NewReadPrivateMessageHandler(notificationReadPrivateMessage))
	router.Handle(notificationWebsocket.ReadRoomMessageTopic, notificationWebsocket.NewReadRoomMessageHandler(notificationReadRoomMessage))

	return router
}

func provideWebsocketServer(upgrader *gorillaWebsocket.Upgrader, manager *websocket.Manager, router *websocket.Router) *websocket.Server {
	return websocket.NewServer(upgrader, manager, router)
}
