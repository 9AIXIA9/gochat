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
	"github.com/ulule/limiter/v3"
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
	websocketLimiter *WebsocketLimiter,
	validator websocket.Validator,
	chatSendPrivateMessage chatApp.SendPrivateMessageUseCase,
	chatSendRoomMessage chatApp.SendRoomMessageUseCase,
	chatConfirmPrivateMessages chatApp.ConfirmPrivateMessagesUseCase,
	chatConfirmRoomMessages chatApp.ConfirmRoomMessagesUseCase,

) *websocket.Router {
	router := websocket.NewRouter(validator)

	middlewares := []websocket.Middleware{
		middleware.NewRateLimitMiddleware((*limiter.Limiter)(websocketLimiter)),
		middleware.NewTimeoutMiddleware(appConfig.Timeout),
		middleware.NewRecoverMiddleware(),
		middleware.NewLoggerMiddleware(),
	}
	if appConfig.OTEL != nil && appConfig.OTEL.Enabled {
		middlewares = append([]websocket.Middleware{middleware.NewTraceMiddleware(appConfig.Name + ".websocket")}, middlewares...)
	}
	router.Use(middlewares...)

	if appConfig.Breaker != nil {
		appConfig.Breaker.Name = appConfig.Name + "_websocket_circuit_breaker"
		router.Use(middleware.NewCircuitBreakMiddleware(appConfig.Breaker))
	}

	router.NoRoute(websocketDelivery.NewNotFoundHandler())

	{
		router.Handle(chatWebsocket.SendPrivateMessageTopic, chatWebsocket.NewSendPrivateMessageHandler(chatSendPrivateMessage, validator))
		router.Handle(chatWebsocket.SendRoomMessageTopic, chatWebsocket.NewSendRoomMessageHandler(chatSendRoomMessage, validator))
		router.Handle(chatWebsocket.ConfirmPrivateMessagesTopic, chatWebsocket.NewConfirmPrivateMessagesHandler(chatConfirmPrivateMessages, validator))
		router.Handle(chatWebsocket.ConfirmRoomMessagesTopic, chatWebsocket.NewConfirmRoomMessagesHandler(chatConfirmRoomMessages, validator))
	}
	return router
}
