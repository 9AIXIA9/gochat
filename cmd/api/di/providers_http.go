package di

import (
	"context"
	"fmt"
	"gochat/docs"
	"gochat/internal/application"
	authApp "gochat/internal/authorization/application"
	chatApp "gochat/internal/chat/application"
	friendshipApp "gochat/internal/friendship/application"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/websocket"
	notificationApp "gochat/internal/notification/application"
	profileApp "gochat/internal/profile/application"
	roomshipApp "gochat/internal/roomship/application"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"gochat/config"
	authHTTP "gochat/internal/authorization/port/http"
	chatHTTP "gochat/internal/chat/port/http"
	"gochat/internal/delivery/http/handler"
	"gochat/internal/delivery/http/middleware"
	friendshipHTTP "gochat/internal/friendship/port/http"
	ginInfra "gochat/internal/infrastructure/gin"
	notificationHTTP "gochat/internal/notification/port/http"
	profileHTTP "gochat/internal/profile/port/http"
	roomshipHTTP "gochat/internal/roomship/port/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
)

var nonBusinessPaths = []string{
	"/healthz",
	"/readyz",
	"/swagger/*any",
	"/metrics",
}

var HTTPSet = wire.NewSet(
	provideHttpRouter,
	provideHttpServer,
	provideWebsocketHandler,
	provideIsReadyChecker,
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
	readPrivateMessages chatApp.ReadPrivateMessagesUseCase,
	readRoomMessages chatApp.ReadRoomMessagesUseCase,
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
	isReady func() bool,
) *gin.Engine {
	// 设置全局环境变量
	ginInfra.SetGlobalEnv(appConfig.Env)

	// 初始化Gin路由器
	router := gin.New()

	// 链路追踪中间件）
	router.Use(middleware.SkipMiddleware(nonBusinessPaths, otelgin.Middleware(appConfig.Name)))

	// swagger base path 保持与路由前缀一致
	docs.SwaggerInfo.BasePath = "/api/v1"

	// 全局中间件
	router.Use(
		middleware.NewRecoverMiddleware(),
		middleware.SkipMiddleware(nonBusinessPaths, middleware.NewLoggerMiddleware()),
		middleware.NewCORSMiddleware(appConfig.CORS),
	)

	// 非业务路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(handler.NewNotFoundHandler())
	router.Any("/healthz", handler.NewHealthCheckHandler())
	router.Any("/readyz", handler.NewReadyCheckHandler(isReady))

	// 业务路由
	baseGroup := router.Group("/api/v1")
	baseGroup.Use(
		middleware.NewRateLimitMiddleware(redisClient, appConfig.RateLimit),
	)

	// 断路器中间件
	if appConfig.Breaker != nil {
		zap.L().Info("Enable Circuit Breaker Middleware")
		appConfig.Breaker.Name = appConfig.Name + "_http_circuit_breaker"
		baseGroup.Use(middleware.NewCircuitBreakMiddleware(appConfig.Breaker))
	}

	authorizationMiddleware := authHTTP.NewAuthorizationMiddleware(parseAccessToken)

	// WebSocket路由
	websocketGroup := baseGroup.Group("/ws")
	websocketGroup.Use(authorizationMiddleware)
	{
		websocketGroup.GET("/", gin.WrapH(websocketHandler))
	}

	baseGroup.Use(
		middleware.NewTimeoutMiddleware(appConfig.Timeout),
	)

	// 授权相关路由（RESTful）
	authGroup := baseGroup.Group("/auth")
	{
		authGroup.POST("/sign-up", authHTTP.NewSignUpHandler(signUp, validator))
		authGroup.POST("/login", authHTTP.NewLoginHandler(login, validator, appConfig.Cookie))
		authGroup.PUT("/tokens/refresh", authHTTP.NewRefreshAccessTokenHandler(refreshAccessToken, validator, appConfig.Cookie))
	}

	// 资料相关路由（RESTful）
	profilesPublic := baseGroup.Group("/profiles")
	{
		profilesPublic.GET("/users/:user_id", profileHTTP.NewGetUserProfileHandler(getUserProfile, validator))
		profilesPublic.GET("/rooms/:room_id", profileHTTP.NewGetRoomProfileHandler(getRoomProfile, validator))
	}

	profilesAuth := baseGroup.Group("/profiles")
	profilesAuth.Use(authorizationMiddleware)
	{
		profilesAuth.GET("/me", profileHTTP.NewGetMyProfileHandler(getUserProfile, validator))
		profilesAuth.PUT("/me", profileHTTP.NewUpdateUserProfileHandler(updateUserProfile, validator))
		// Prefer path param; keep old body-based route for compatibility below
		profilesAuth.PUT("/rooms/:room_id", profileHTTP.NewUpdateRoomProfileHandler(updateRoomProfile, validator))
	}

	// 聊天相关路由（RESTful）
	chatsGroup := baseGroup.Group("/chats")
	chatsGroup.Use(authorizationMiddleware)
	{
		chatsGroup.GET("/private-messages/:user_id", chatHTTP.NewListPrivateMessagesHandler(listPrivateMessages, validator))
		chatsGroup.GET("/rooms/messages/:room_id", chatHTTP.NewListRoomMessagesHandler(listRoomMessages, validator))
		chatsGroup.POST("/private-messages", chatHTTP.NewSendPrivateMessageHandler(sendPrivateMessage, validator))
		chatsGroup.POST("/rooms/messages", chatHTTP.NewSendRoomMessageHandler(sendRoomMessage, validator))
		chatsGroup.PUT("/private-messages/read", chatHTTP.NewReadPrivateMessagesHandler(readPrivateMessages, validator))
		chatsGroup.PUT("/rooms/messages/read", chatHTTP.NewReadRoomMessagesHandler(readRoomMessages, validator))
	}

	// 房间与成员相关路由（RESTful）
	// 公共：获取房间成员列表
	roomsPublic := baseGroup.Group("/rooms")
	{
		roomsPublic.GET("/:room_id/members", roomshipHTTP.NewListRoomMembersHandler(listRoomMembers, validator))
	}

	// 授权：管理我的房间与成员关系
	roomsAuth := baseGroup.Group("/rooms")
	roomsAuth.Use(authorizationMiddleware)
	{
		// Rooms
		roomsAuth.POST("/", roomshipHTTP.NewCreateRoomHandler(createRoom, validator))
		// 我加入的房间列表
		roomsAuth.GET("/", roomshipHTTP.NewListRoomshipsHandler(listRoomships, validator))
		// 退出房间
		roomsAuth.DELETE("/:room_id/members/me", roomshipHTTP.NewLeaveRoomHandler(leaveRoom, validator))
	}

	// 成员请求
	roomRequests := baseGroup.Group("/rooms/requests")
	roomRequests.Use(authorizationMiddleware)
	{
		roomRequests.GET("/", roomshipHTTP.NewListMemberRequestsHandler(listMemberRequests, validator))
		roomRequests.POST("/", roomshipHTTP.NewSendMemberRequestHandler(sendMemberRequest, validator))
		roomRequests.PUT("/:request_id/agree", roomshipHTTP.NewAgreeMemberRequestHandler(agreeMemberRequest, validator))
		roomRequests.PUT("/:request_id/refuse", roomshipHTTP.NewRefuseMemberRequestHandler(refuseMemberRequest, validator))
	}

	// 好友功能路由（RESTful）
	friendshipsGroup := baseGroup.Group("/friendships")
	friendshipsGroup.Use(authorizationMiddleware)
	{
		friendshipsGroup.GET("/", friendshipHTTP.NewListFriendshipsHandler(listFriendships, validator))
	}

	friendshipRequestsGroup := baseGroup.Group("/friendship-requests")
	friendshipRequestsGroup.Use(authorizationMiddleware)
	{
		friendshipRequestsGroup.GET("/", friendshipHTTP.NewListFriendRequestsHandler(listFriendRequests, validator))
		friendshipRequestsGroup.POST("/", friendshipHTTP.NewSendFriendRequestHandler(sendFriendRequest, validator))
		friendshipRequestsGroup.PUT("/:request_id/agree", friendshipHTTP.NewAgreeFriendRequestHandler(agreeFriendRequest, validator))
		friendshipRequestsGroup.PUT("/:request_id/refuse", friendshipHTTP.NewRefuseFriendRequestHandler(refuseFriendRequest, validator))
	}

	// 通知功能路由（RESTful）
	notificationsGroup := baseGroup.Group("/notifications")
	notificationsGroup.Use(authorizationMiddleware)
	{
		notificationsGroup.GET("/system-messages", notificationHTTP.NewListSystemMessagesHandler(listSystemMessages, validator))
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

func provideIsReadyChecker(
	db *gorm.DB,
	rdb *redis.Client,
	producer *kafka.Producer,
	consumers []*kafkaInfra.Consumer,
) func() bool {
	return func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			dbConn, err := db.DB()
			if err != nil {
				return err
			}
			return dbConn.PingContext(ctx)
		})

		group.Go(func() error {
			return rdb.Ping(ctx).Err()
		})

		group.Go(func() error {
			// Flush returns remaining buffered messages; treat leftover as unhealthy
			if remaining := producer.Flush(1000); remaining > 0 {
				return fmt.Errorf("kafka producer has %d pending messages", remaining)
			}
			return nil
		})

		for _, c := range consumers {
			consumer := c
			group.Go(func() error {
				return consumer.Ping(ctx)
			})
		}

		return group.Wait() == nil
	}
}
