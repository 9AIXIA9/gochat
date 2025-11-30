package di

import (
	"fmt"
	authApp "gochat/internal/authorization/application"
	chatApp "gochat/internal/chat/application"
	roomshipApp "gochat/internal/roomship/application"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"gochat/config"
	authHTTP "gochat/internal/authorization/port/http"
	chatHTTP "gochat/internal/chat/port/http"
	"gochat/internal/delivery/http/handler"
	"gochat/internal/delivery/http/middleware"
	ginInfra "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/infrastructure/validator"
	"gochat/internal/infrastructure/websocket"
	roomshipHTTP "gochat/internal/roomship/port/http"
)

var HTTPSet = wire.NewSet(
	provideHttpRouter,
	provideHttpServer,
)

func provideHttpRouter(
	appConfig *config.App,
	signUp authApp.SignUpUseCase,
	login authApp.LoginUseCase,
	refreshAccessToken authApp.RefreshAccessTokenUseCase,
	parseAccessToken authApp.ParseAccessTokenUseCase,
	sendPrivateMessage chatApp.SendPrivateMessageUseCase,
	sendRoomMessage chatApp.SendRoomMessageUseCase,
	createRoom roomshipApp.CreateRoomUseCase,
	joinRoom roomshipApp.JoinRoomUseCase,
	leaveRoom roomshipApp.LeaveRoomUseCase,
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

	router.NoRoute(handler.NewNotFoundHandler()) // 404处理器

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
		authorizationGroup.POST("/sign_up", authHTTP.NewSignUpHandler(signUp, validator))
		authorizationGroup.POST("/login", authHTTP.NewLoginHandler(login, validator, appConfig.Cookie))
		authorizationGroup.GET("/refresh_access_token", authHTTP.NewRefreshAccessTokenHandler(refreshAccessToken, validator, appConfig.Cookie))
	}

	authorizationMiddleware := authHTTP.NewAuthorizationMiddleware(parseAccessToken)

	// 聊天相关路由
	chatGroup := baseGroup.Group("/chat")
	chatGroup.Use(authorizationMiddleware)
	{
		chatGroup.POST("/private", chatHTTP.NewSendPrivateMessageHandler(sendPrivateMessage, validator))
		chatGroup.POST("/room", chatHTTP.NewSendRoomMessageHandler(sendRoomMessage, validator))
	}

	// 房间功能路由
	roomshipGroup := baseGroup.Group("/roomship")
	roomshipGroup.Use(authorizationMiddleware)
	{
		roomshipGroup.POST("/room", roomshipHTTP.NewCreateRoomHandler(createRoom, validator))
		roomshipGroup.POST("/room/member", roomshipHTTP.NewJoinRoomHandler(joinRoom, validator))
		roomshipGroup.DELETE("/room/member", roomshipHTTP.NewLeaveRoomHandler(leaveRoom, validator))
	}

	// WebSocket路由
	websocketGroup := baseGroup.Group("/ws")
	websocketGroup.Use(authorizationMiddleware)
	{
		websocketGroup.GET("/", gin.WrapH(websocketServer))
	}

	return router
}

func provideHttpServer(appConfig *config.App, router *gin.Engine) *ginInfra.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port),
		Handler: router,
	}
	return ginInfra.NewServer(router, srv)
}
