package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/logger"
	"gochat/internal/middleware"
	"gochat/internal/presentation"
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
		public.POST("/signup", presentation.SignupHandlerFunc(deps.SignupUsecase))
		public.POST("/login", presentation.LoginHandlerFunc(deps.LoginUsecase))
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(deps.AuthUsecase))
	{
		// 房间相关
		protected.POST("/rooms", presentation.CreateRoomHandlerFunc(deps.CreateRoomUsecase))
		protected.POST("/rooms/:number/join", presentation.JoinRoomHandlerFunc(deps.JoinRoomUsecase))
		protected.DELETE("/rooms/:number/leave", presentation.LeaveRoomHandlerFunc(deps.LeaveRoomUsecase))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		presentation.ResponseSuccess(c, domain.NewSuccessResponse("服务器正常工作"))
	})
}
