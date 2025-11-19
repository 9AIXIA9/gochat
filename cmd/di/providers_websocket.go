package di

import (
	"gochat/config"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/infrastructure/websocket"
	notificationUsecase "gochat/internal/notification/application/usecase"
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
	notificationMessageRead notificationUsecase.MessageReadUseCase,
) *websocket.Router {
	router := websocket.NewRouter()

	router.Handle(notificationWebsocket.MessageReadRequestTopic, notificationWebsocket.NewMessageReadHandler(notificationMessageRead))

	return router
}

func provideWebsocketServer(manager *websocket.Manager, router *websocket.Router, eventRepo *repository.EventRepository, eventIDGen event.IDGenerator) *websocket.Server {
	return websocket.NewServer(manager, router, eventRepo, eventIDGen)
}
