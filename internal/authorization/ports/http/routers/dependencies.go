package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/config"
	"gochat/internal/shared/http/middlewares"
	ginutils "gochat/internal/shared/infrastructure/gin"
)

type Dependencies struct {
	AuthorizationMiddleware   gin.HandlerFunc
	SignUpUseCase             domain.SignUpUseCase
	LoginUseCase              domain.LoginUseCase
	RefreshAccessTokenUseCase domain.RefreshAccessTokenUseCase
	Validator                 *ginutils.Validator
	RedisClient               *redis.Client
	RateLimitConfig           *middlewares.RateLimitConfig
	CookieConfig              *config.Cookie
}
