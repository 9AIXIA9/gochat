package di

import (
	"gochat/config"
	websocketDelivery "gochat/internal/delivery/websocket"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/infrastructure/websocket"
	notificationUsecase "gochat/internal/notification/application"
	notificationWebsocket "gochat/internal/notification/port/websocket"
	"gochat/internal/shared/event"

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

func provideWebsocketManager(upgrader *gorillaWebsocket.Upgrader, metrics *prometheus.Metrics) *websocket.Manager {
	m := websocket.NewManager(upgrader)
	m.SetMetrics(metrics)
	return m
}

func provideWebsocketRouter(
	appConfig *config.App,
	notificationPrivateMessageRead notificationUsecase.PrivateMessageReadUseCase,
	notificationRoomMessageRead notificationUsecase.RoomMessageReadUseCase,
) *websocket.Router {
	router := websocket.NewRouter()

	if appConfig.Telemetry != nil && appConfig.Telemetry.Enabled && appConfig.Telemetry.TraceEnabled {
		router.Use(websocketDelivery.NewTelemetryMiddleware(appConfig.Name))
	}

	router.Handle(notificationWebsocket.PrivateMessageReadRequestTopic, notificationWebsocket.NewPrivateMessageReadHandler(notificationPrivateMessageRead))
	router.Handle(notificationWebsocket.RoomMessageReadRequestTopic, notificationWebsocket.NewRoomMessageReadHandler(notificationRoomMessageRead))

	return router
}

func provideWebsocketServer(manager *websocket.Manager, router *websocket.Router, eventRepo *repository.EventRepository, eventIDGen event.IDGenerator) *websocket.Server {
	return websocket.NewServer(manager, router, eventRepo, eventIDGen)
}
