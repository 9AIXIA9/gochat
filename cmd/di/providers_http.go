package di

import (
	"fmt"
	"gochat/docs"
	"gochat/internal/application"
	authApp "gochat/internal/authorization/application"
	chatApp "gochat/internal/chat/application"
	friendshipApp "gochat/internal/friendship/application"
	"gochat/internal/infrastructure/websocket"
	notificationApp "gochat/internal/notification/application"
	profileApp "gochat/internal/profile/application"
	roomshipApp "gochat/internal/roomship/application"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"gochat/config"
	authHTTP "gochat/internal/authorization/port/http"
	chatHTTP "gochat/internal/chat/port/http"
	"gochat/internal/delivery/http/handler"
	"gochat/internal/delivery/http/middleware"
	friendshipHTTP "gochat/internal/friendship/port/http"
	ginInfra "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/prometheus"
	notificationHTTP "gochat/internal/notification/port/http"
	profileHTTP "gochat/internal/profile/port/http"
	roomshipHTTP "gochat/internal/roomship/port/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var HTTPSet = wire.NewSet(
	provideHttpRouter,
	provideHttpServer,
	provideWebsocketHandler,
)

func provideHttpRouter(
	appConfig *config.App,
	signUp authApp.SignUpUseCase,
	login authApp.LoginUseCase,
	refreshAccessToken authApp.RefreshAccessTokenUseCase,
	parseAccessToken authApp.ParseAccessTokenUseCase,
	getUserProfile profileApp.GetUserProfileUseCase,
	getRoomProfile profileApp.GetRoomProfileUseCase,
	updateUserProfile profileApp.UpdateUserProfileUseCase,
	updateRoomProfile profileApp.UpdateRoomProfileUseCase,
	sendPrivateMessage chatApp.SendPrivateMessageUseCase,
	sendRoomMessage chatApp.SendRoomMessageUseCase,
	listPrivateMessages chatApp.ListPrivateMessagesUseCase,
	listRoomMessages chatApp.ListRoomMessagesUseCase,
	createRoom roomshipApp.CreateRoomUseCase,
	listRoomMembers roomshipApp.ListRoomMembersUseCase,
	sendMemberRequest roomshipApp.SendMemberRequestUseCase,
	agreeMemberRequest roomshipApp.AgreeMemberRequestUseCase,
	refuseMemberRequest roomshipApp.RefuseMemberRequestUseCase,
	listMemberRequests roomshipApp.ListMemberRequestsUseCase,
	listRoomships roomshipApp.ListRoomshipsUseCase,
	leaveRoom roomshipApp.LeaveRoomUseCase,
	sendFriendRequest friendshipApp.SendFriendRequestUseCase,
	agreeFriendRequest friendshipApp.AgreeFriendRequestUseCase,
	refuseFriendRequest friendshipApp.RefuseFriendRequestUseCase,
	listFriendships friendshipApp.ListFriendshipsUseCase,
	listFriendRequests friendshipApp.ListFriendRequestsUseCase,
	listSystemMessages notificationApp.ListSystemMessagesUseCase,
	validator ginInfra.Validator,
	redisClient *redis.Client,
	websocketHandler *handler.WebsocketHandler,
	metrics *prometheus.Metrics,
) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(appConfig.Env)

	// 初始化Gin路由器
	router := gin.New()

	// swagger base path 保持与路由前缀一致
	docs.SwaggerInfo.BasePath = "/api/v1"

	// 全局中间件
	router.Use(
		middleware.NewRecoverMiddleware(),
		middleware.NewLoggerMiddleware(),
		middleware.NewCORSMiddleware(appConfig.CORS),
	)

	// 非业务路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(handler.NewNotFoundHandler())
	router.Any("/metrics", gin.WrapH(metrics.Handler()))
	router.Any("/health_check", handler.NewHealthCheckHandler())

	// 业务路由
	baseGroup := router.Group("/api/v1")
	baseGroup.Use(
		middleware.NewRateLimitMiddleware(redisClient, appConfig.RateLimit),
		metrics.GinMiddleware(),
	)

	authorizationMiddleware := authHTTP.NewAuthorizationMiddleware(parseAccessToken)

	// 遥测追踪中间件
	if appConfig.Telemetry != nil && appConfig.Telemetry.Enabled && appConfig.Telemetry.TraceEnabled {
		baseGroup.Use(middleware.NewTelemetryMiddleware(appConfig.Name))
	}

	// 断路器中间件
	if appConfig.Breaker != nil {
		zap.L().Info("Enable Circuit Breaker Middleware")
		appConfig.Breaker.Name = appConfig.Name + "_http_circuit_breaker"
		baseGroup.Use(middleware.NewCircuitBreakMiddleware(appConfig.Breaker))
	}

	// WebSocket路由
	websocketGroup := baseGroup.Group("/ws")
	websocketGroup.Use(authorizationMiddleware)
	{
		websocketGroup.GET("/", gin.WrapH(websocketHandler))
	}

	baseGroup.Use(
		middleware.NewTimeoutMiddleware(appConfig.Timeout),
	)

	// 授权相关路由
	authorizationGroup := baseGroup.Group("/authorization")
	{
		authorizationGroup.POST("/sign_up", authHTTP.NewSignUpHandler(signUp, validator))
		authorizationGroup.POST("/login", authHTTP.NewLoginHandler(login, validator, appConfig.Cookie))
		authorizationGroup.GET("/refresh_access_token", authHTTP.NewRefreshAccessTokenHandler(refreshAccessToken, validator, appConfig.Cookie))
	}

	// 聊天相关路由
	profileGroup := baseGroup.Group("/profile")
	{
		profileGroup.GET("/user/:user_id", profileHTTP.NewGetUserProfileHandler(getUserProfile, validator))
		profileGroup.GET("/room/:room_id", profileHTTP.NewGetRoomProfileHandler(getRoomProfile, validator))
	}

	profileGroup.Use(authorizationMiddleware)
	{
		profileGroup.PUT("/me", profileHTTP.NewUpdateUserProfileHandler(updateUserProfile, validator))
		profileGroup.PUT("/room", profileHTTP.NewUpdateRoomProfileHandler(updateRoomProfile, validator))
	}

	// 聊天相关路由
	chatGroup := baseGroup.Group("/chat")
	chatGroup.Use(authorizationMiddleware)
	{
		chatGroup.GET("/private", chatHTTP.NewListPrivateMessagesHandler(listPrivateMessages, validator))
		chatGroup.GET("/room", chatHTTP.NewListRoomMessagesHandler(listRoomMessages, validator))
		chatGroup.POST("/private", chatHTTP.NewSendPrivateMessageHandler(sendPrivateMessage, validator))
		chatGroup.POST("/room", chatHTTP.NewSendRoomMessageHandler(sendRoomMessage, validator))
	}

	// 房间功能路由
	roomshipGroup := baseGroup.Group("/roomship")
	{
		roomshipGroup.GET("/room/:room_id", roomshipHTTP.NewListRoomMembersHandler(listRoomMembers, validator))
	}

	roomshipGroup.Use(authorizationMiddleware)
	{
		// Rooms
		roomshipGroup.POST("/room", roomshipHTTP.NewCreateRoomHandler(createRoom, validator))

		// Roomship
		roomshipGroup.GET("/", roomshipHTTP.NewListRoomshipsHandler(listRoomships, validator))
		roomshipGroup.DELETE("/room/:room_id", roomshipHTTP.NewLeaveRoomHandler(leaveRoom, validator))

		// Member Requests
		roomshipGroup.GET("/request", roomshipHTTP.NewListMemberRequestsHandler(listMemberRequests, validator))
		roomshipGroup.POST("/request", roomshipHTTP.NewSendMemberRequestHandler(sendMemberRequest, validator))
		roomshipGroup.PUT("/request/:request_id/agree", roomshipHTTP.NewAgreeMemberRequestHandler(agreeMemberRequest, validator))
		roomshipGroup.PUT("/request/:request_id/refuse", roomshipHTTP.NewRefuseMemberRequestHandler(refuseMemberRequest, validator))
	}

	// 好友功能路由
	friendshipGroup := baseGroup.Group("/friendship")
	friendshipGroup.Use(authorizationMiddleware)
	{
		friendshipGroup.GET("/", friendshipHTTP.NewListFriendshipsHandler(listFriendships, validator))
		friendshipGroup.GET("/request", friendshipHTTP.NewListFriendRequestsHandler(listFriendRequests, validator))
		friendshipGroup.POST("/request", friendshipHTTP.NewSendFriendRequestHandler(sendFriendRequest, validator))
		friendshipGroup.PUT("/request/:request_id/agree", friendshipHTTP.NewAgreeFriendRequestHandler(agreeFriendRequest, validator))
		friendshipGroup.PUT("/request/:request_id/refuse", friendshipHTTP.NewRefuseFriendRequestHandler(refuseFriendRequest, validator))
	}

	// 通知功能路由
	notificationGroup := baseGroup.Group("/notification")
	notificationGroup.Use(authorizationMiddleware)
	{
		notificationGroup.GET("/system", notificationHTTP.NewListSystemMessagesHandler(listSystemMessages, validator))
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

func provideWebsocketHandler(
	upgrader *gorillaWebsocket.Upgrader,
	manager *websocket.Manager,
	router *websocket.Router,
	userSessionStartedUseCase application.UserSessionStartedUseCase,
) *handler.WebsocketHandler {
	return handler.NewWebsocketHandler(
		upgrader,
		manager,
		router,
		userSessionStartedUseCase,
	)
}
