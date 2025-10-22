package routers

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/authorization/ports/http/handlers"
	"gochat/internal/shared/http"
	"gochat/internal/shared/http/middlewares"
	ginutils "gochat/internal/shared/infrastructure/gin"
)

func Setup(group *gin.RouterGroup, dependencies *Dependencies) {
	group.Use(middlewares.RateLimit(dependencies.RedisClient, dependencies.RateLimitConfig))
	{
		group.POST("/sign_up", handlers.NewSignUp(dependencies.SignUpUseCase, dependencies.Validator))
		group.POST("/login", handlers.NewLogin(dependencies.LoginUseCase, dependencies.Validator, dependencies.CookieConfig))
		group.GET("/refresh_access_token", handlers.NewRefreshAccessToken(dependencies.RefreshAccessTokenUseCase, dependencies.Validator, dependencies.CookieConfig))
	}

	testGroup := group.Group("/")
	testGroup.Use(dependencies.AuthorizationMiddleware)
	testGroup.GET("test", func(ginContext *gin.Context) {
		ginutils.Response(ginContext, http.ResponseSuccess)
	})
}
