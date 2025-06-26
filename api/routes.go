package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
	"gochat/internal/middleware"
	"gochat/internal/presentation"
)

// Dependencies 依赖注入结构体
type Dependencies struct {
	Logger            domain.Logger
	AuthUsecase       domain.AuthUsecase
	SignupUsecase     domain.SignupUsecase
	LoginUsecase      domain.LoginUsecase
	CreateRoomUsecase domain.CreateRoomUsecase
	JoinRoomUsecase   domain.JoinRoomUsecase
	ExitRoomUsecase   domain.ExitRoomUsecase
}

func Setup(deps *Dependencies) *gin.Engine {
	// 创建gin引擎
	r := gin.Default()

	// 注册全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(deps.Logger))
	r.Use(middleware.Error(deps.Logger))

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
		public.POST("/signup", presentation.SignupHandlerFunc(deps.SignupUsecase, deps.Logger))
		public.POST("/login", presentation.LoginHandlerFunc(deps.LoginUsecase, deps.Logger))
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(deps.AuthUsecase))
	{
		// 房间相关
		protected.POST("/rooms", presentation.CreateRoomHandlerFunc(deps.CreateRoomUsecase, deps.Logger))
		protected.POST("/rooms/:number/join", presentation.JoinRoomHandlerFunc(deps.JoinRoomUsecase, deps.Logger))
		protected.DELETE("/rooms/:number/exit", presentation.ExitRoomHandlerFunc(deps.ExitRoomUsecase, deps.Logger))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		presentation.ResponseSuccess(c, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})
}
