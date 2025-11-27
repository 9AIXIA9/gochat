package di

import (
	"gochat/config"
	websocketDelivery "gochat/internal/delivery/websocket"
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
	notificationPrivateMessageRead notificationApp.PrivateMessageReadUseCase,
	notificationRoomMessageRead notificationApp.RoomMessageReadUseCase,
) *websocket.Router {
	router := websocket.NewRouter(websocketDelivery.NewNotFoundHandler())

	router.Handle(notificationWebsocket.PrivateMessageReadTopic, notificationWebsocket.NewPrivateMessageReadHandler(notificationPrivateMessageRead))
	router.Handle(notificationWebsocket.RoomMessageReadTopic, notificationWebsocket.NewRoomMessageReadHandler(notificationRoomMessageRead))

	return router
}

func provideWebsocketServer(upgrader *gorillaWebsocket.Upgrader, manager *websocket.Manager, router *websocket.Router) *websocket.Server {
	return websocket.NewServer(upgrader, manager, router)
}
