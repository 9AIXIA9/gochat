package api

import (
	"gochat/config"
	"gochat/internal/authorization/application/usecase"
	authorizationHttp "gochat/internal/authorization/port/http"
	"gochat/internal/delivery/http/middleware"
	ginutils "gochat/internal/infrastructure/validator"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	signUpUseCase usecase.SignUpUseCase,
	loginUseCase usecase.LoginUseCase,
	refreshAccessTokenUseCase usecase.RefreshAccessTokenUseCase,
	parseAccessTokenUseCase usecase.ParseAccessTokenUseCase,
	validator *ginutils.Validator,
	redisClient *redis.Client,
	rateLimitConfig *middleware.RateLimitConfig,
	cookieConfig *config.Cookie,
	CORSConfig *middleware.CORSConfig,
) *gin.Engine {
	engine := gin.New()

	engine.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewCORSMiddleware(CORSConfig),
		middleware.NewRateLimitMiddleware(redisClient, rateLimitConfig),
	)

	baseGroup := engine.Group("/api/v1")

	authorizationGroup := baseGroup.Group("/authorization")

	authorizationGroup.Use()
	{
		authorizationGroup.POST("/sign_up", authorizationHttp.NewSignUpHandler(signUpUseCase, validator))
		authorizationGroup.POST("/login", authorizationHttp.NewLoginHandler(loginUseCase, validator, cookieConfig))
		authorizationGroup.GET("/refresh_access_token", authorizationHttp.NewRefreshAccessTokenHandler(refreshAccessTokenUseCase, validator, cookieConfig))
	}

	//authorizationMiddleware := authorizationHttp.NewAuthorizationMiddleware(parseAccessTokenUseCase)
	_ = authorizationHttp.NewAuthorizationMiddleware(parseAccessTokenUseCase)

	return engine
}
