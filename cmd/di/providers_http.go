package di

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"gochat/config"
	authorizationUsecase "gochat/internal/authorization/application/usecase"
	authorizationHttp "gochat/internal/authorization/port/http"
	chatUsecase "gochat/internal/chat/application/usecase"
	chatHttp "gochat/internal/chat/port/http"
	"gochat/internal/delivery/http/handler"
	"gochat/internal/delivery/http/middleware"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/infrastructure/validator"
	"gochat/internal/infrastructure/websocket"
	socialUseCase "gochat/internal/social/application/usecase"
	socialHttp "gochat/internal/social/port/http"
)

var HTTPSet = wire.NewSet(
	provideHttpRouter,
	provideHttpServer,
)

func provideHttpRouter(
	appConfig *config.App,
	signUp authorizationUsecase.SignUpUseCase,
	login authorizationUsecase.LoginUseCase,
	refreshAccessToken authorizationUsecase.RefreshAccessTokenUseCase,
	parseAccessToken authorizationUsecase.ParseAccessTokenUseCase,
	sendPrivateMessage chatUsecase.SendPrivateMessageUseCase,
	sendRoomMessage chatUsecase.SendRoomMessageUseCase,
	createRoom socialUseCase.CreateRoomUseCase,
	joinRoom socialUseCase.JoinRoomUseCase,
	leaveRoom socialUseCase.LeaveRoomUseCase,
	validator *validator.Validator,
	redisClient *redis.Client,
	websocketServer *websocket.Server,
	metrics *prometheus.Metrics,
) *gin.Engine {

	// 初始化Gin路由器
	router := gin.New()

	// 全局中间件栈
	router.Use(
		middleware.NewRecoverMiddleware(),            // 恢复中间件
		middleware.NewLoggerMiddleware(),             // 日志中间件
		middleware.NewCORSMiddleware(appConfig.CORS), // CORS中间件
	)

	// 暴露Prometheus指标端点
	router.Any("/metrics", gin.WrapH(metrics.Handler()))

	// 健康检查端点
	router.Any("/health_check", handler.NewHealthCheckHandler())

	// API路由分组 - 基础路径
	baseGroup := router.Group("/api/v1")

	baseGroup.Use(
		middleware.NewRateLimitMiddleware(redisClient, appConfig.RateLimit), // 限流中间件
		metrics.GinMiddleware(), // Prometheus指标中间件
	)

	// 遥测追踪中间件（如果启用）
	if appConfig.Telemetry != nil && appConfig.Telemetry.Enabled && appConfig.Telemetry.TraceEnabled {
		baseGroup.Use(middleware.NewTelemetryMiddleware(appConfig.Name))
	}

	// 授权相关路由
	authorizationGroup := baseGroup.Group("/authorization")
	authorizationGroup.Use()
	{
		authorizationGroup.POST("/sign_up", authorizationHttp.NewSignUpHandler(signUp, validator))
		authorizationGroup.POST("/login", authorizationHttp.NewLoginHandler(login, validator, appConfig.Cookie))
		authorizationGroup.GET("/refresh_access_token", authorizationHttp.NewRefreshAccessTokenHandler(refreshAccessToken, validator, appConfig.Cookie))
	}

	authorizationMiddleware := authorizationHttp.NewAuthorizationMiddleware(parseAccessToken)

	// 聊天相关路由
	chatGroup := baseGroup.Group("/chat")
	chatGroup.Use(authorizationMiddleware)
	{
		chatGroup.POST("/private", chatHttp.NewSendPrivateMessageHandler(sendPrivateMessage, validator))
		chatGroup.POST("/room", chatHttp.NewSendRoomMessageHandler(sendRoomMessage, validator))
	}

	// 社交功能路由
	socialGroup := baseGroup.Group("/social")
	socialGroup.Use(authorizationMiddleware)
	{
		socialGroup.POST("/room", socialHttp.NewCreateRoomHandler(createRoom, validator))
		socialGroup.POST("/room/member", socialHttp.NewJoinRoomHandler(joinRoom, validator))
		socialGroup.DELETE("/room/member", socialHttp.NewLeaveRoomHandler(leaveRoom, validator))
	}

	// WebSocket路由
	websocketGroup := baseGroup.Group("/ws")
	websocketGroup.Use(authorizationMiddleware)
	{
		websocketGroup.GET("/", handler.NewWebsocketHandler(websocketServer))
	}

	return router
}

func provideHttpServer(appConfig *config.App, router *gin.Engine) *ginutils.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port),
		Handler: router,
	}
	return ginutils.NewServer(router, srv)
}
