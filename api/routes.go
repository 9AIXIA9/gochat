package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/handler"
	"gochat/internal/infra/logger"
	"gochat/internal/middleware"
)

// Dependencies 依赖注入结构体
type Dependencies struct {
	Config            *config.Config
	AuthUsecase       domain.AuthUsecase
	SignupUsecase     domain.SignupUsecase
	LoginUsecase      domain.LoginUsecase
	CreateRoomUsecase domain.CreateRoomUsecase
	JoinRoomUsecase   domain.JoinRoomUsecase
	LeaveRoomUsecase  domain.LeaveRoomUsecase
}

func Setup(deps *Dependencies) *gin.Engine {
	// 创建gin引擎
	r := gin.New(logger.GinOption())

	// 注册全局中间件
	r.Use(middleware.CORS())

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
		// 用户相关
		public.POST("/signup", handler.Adapter[handler.SignupRequest, domain.SignupRequest](deps.SignupUsecase))
		public.POST("/login", handler.Adapter[handler.LoginRequest, domain.LoginRequest](deps.LoginUsecase))
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(deps.AuthUsecase))
	{
		// 房间相关
		protected.POST("/rooms", handler.Adapter[handler.CreateRoomRequest, domain.CreateRoomRequest](deps.CreateRoomUsecase))
		protected.POST("/rooms/:number/join", handler.Adapter[handler.JoinRoomRequest, domain.JoinRoomRequest](deps.JoinRoomUsecase))
		protected.DELETE("/rooms/:number/leave", handler.Adapter[handler.LeaveRoomRequest, domain.LeaveRoomRequest](deps.LeaveRoomUsecase))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		handler.ResponseSuccess(c, domain.NewSuccessResponse("服务器正常工作"))
	})
}
