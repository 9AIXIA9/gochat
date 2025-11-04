package api

import (
	"gochat/config"
	authorizationUseCase "gochat/internal/authorization/application/usecase"
	authorizationHttp "gochat/internal/authorization/port/http"
	chatUseCase "gochat/internal/chat/application/usecase"
	chatHttp "gochat/internal/chat/port/http"
	"gochat/internal/delivery/http/handler"
	"gochat/internal/delivery/http/middleware"
	ginutils "gochat/internal/infrastructure/validator"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	signUpUseCase authorizationUseCase.SignUpUseCase,
	loginUseCase authorizationUseCase.LoginUseCase,
	refreshAccessTokenUseCase authorizationUseCase.RefreshAccessTokenUseCase,
	parseAccessTokenUseCase authorizationUseCase.ParseAccessTokenUseCase,
	sendPrivateMessageUseCase chatUseCase.SendPrivateMessageUseCase,
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

	engine.Any("/health_check", handler.NewHealthCheckHandler())

	baseGroup := engine.Group("/api/v1")

	authorizationGroup := baseGroup.Group("/authorization")

	authorizationGroup.Use()
	{
		authorizationGroup.POST("/sign_up", authorizationHttp.NewSignUpHandler(signUpUseCase, validator))
		authorizationGroup.POST("/login", authorizationHttp.NewLoginHandler(loginUseCase, validator, cookieConfig))
		authorizationGroup.GET("/refresh_access_token", authorizationHttp.NewRefreshAccessTokenHandler(refreshAccessTokenUseCase, validator, cookieConfig))
	}

	authorizationMiddleware := authorizationHttp.NewAuthorizationMiddleware(parseAccessTokenUseCase)

	chatGroup := baseGroup.Group("/chat")
	chatGroup.Use(authorizationMiddleware)
	{
		chatGroup.POST("/private", chatHttp.NewSendPrivateMessageHandler(sendPrivateMessageUseCase, validator))
	}

	return engine
}
