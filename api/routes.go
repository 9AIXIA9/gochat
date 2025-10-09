package api

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/handler"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/middleware"
)

// Dependencies 依赖注入结构体
type Dependencies struct {
	Config                    *config.Config
	WebsocketManager          *manager.Manager
	RedisClient               *redis.Client
	AuthUsecase               domain.AuthUsecase
	SignupUsecase             domain.SignupUsecase
	LoginUsecase              domain.LoginUsecase
	RefreshTokenUsecase       domain.RefreshTokenUsecase
	CreateRoomUsecase         domain.CreateRoomUsecase
	JoinRoomUsecase           domain.JoinRoomUsecase
	LeaveRoomUsecase          domain.LeaveRoomUsecase
	SendPrivateMessageUsecase domain.SendPrivateMessageUsecase
	SendRoomMessageUsecase    domain.SendRoomMessageUsecase
	UserConnectedUsecase      domain.UserConnectedUsecase
}

func Setup(deps *Dependencies) *gin.Engine {
	// 创建gin引擎
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORS(deps.Config.CORS))
	r.Use(middleware.RateLimit(deps.RedisClient, deps.Config.RateLimit))

	// 注册路由
	setup(r, deps)

	return r
}

func setup(r *gin.Engine, deps *Dependencies) {
	// API路由组
	api := r.Group("/api/v1")

	// 公开路由（不需要认证）
	public := api.Group("/")
	{
		public.POST("/signup", handler.Signup(deps.SignupUsecase, deps.Config.Timeout))
		public.POST("/login", handler.Login(deps.LoginUsecase, deps.Config.Cookie, deps.Config.Token.Refresh, deps.Config.Timeout))
		public.GET("/refresh/token", handler.RefreshToken(deps.RefreshTokenUsecase, deps.Config.Cookie, deps.Config.Token.Refresh, deps.Config.Timeout))
	}

	// 需要认证的路由
	httpsProtected := api.Group("/")
	httpsProtected.Use(middleware.HTTPSJWTAuth(deps.AuthUsecase))
	{
		// 房间相关
		httpsProtected.POST("/rooms", handler.CreateRoom(deps.CreateRoomUsecase, deps.Config.Timeout))
		httpsProtected.POST("/rooms/:number/join", handler.JoinRoom(deps.JoinRoomUsecase, deps.Config.Timeout))
		httpsProtected.DELETE("/rooms/:number/leave", handler.LeaveRoom(deps.LeaveRoomUsecase, deps.Config.Timeout))

		// 消息相关
		httpsProtected.POST("/private/message/:to", handler.SendPrivateMessage(deps.SendPrivateMessageUsecase, deps.Config.Timeout))
		httpsProtected.POST("/room/message/:to", handler.SendRoomMessage(deps.SendRoomMessageUsecase, deps.Config.Timeout))
	}

	// websocket长连接
	wssProtected := api.Group("/")
	wssProtected.Use(middleware.WebsocketJWTAuth(deps.AuthUsecase))
	{
		wssProtected.GET("/ws", handler.Websocket(deps.UserConnectedUsecase, deps.WebsocketManager))
	}
}
