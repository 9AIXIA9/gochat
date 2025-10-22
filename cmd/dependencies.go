package main

import (
	"fmt"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authorizationConverters "gochat/internal/authorization/infrastructure/persistence/converters"
	authorizationModels "gochat/internal/authorization/infrastructure/persistence/models"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	authorizationRouter "gochat/internal/authorization/ports/http/routers"
	authorizationUseCase "gochat/internal/authorization/usecase"
	"gochat/internal/shared/config"
	"gochat/internal/shared/http/middlewares"
	"gochat/internal/shared/infrastructure/bcrypt"
	ginutils "gochat/internal/shared/infrastructure/gin"
	"gochat/internal/shared/infrastructure/gorm"
	"gochat/internal/shared/infrastructure/persistence/converters"
	"gochat/internal/shared/infrastructure/persistence/models"
	"gochat/internal/shared/infrastructure/persistence/repository"
	"gochat/internal/shared/infrastructure/redis"
	"gochat/internal/shared/infrastructure/snowflake"
	"gochat/internal/shared/infrastructure/uuid"
	"gochat/internal/shared/infrastructure/zap"
)

type Dependencies struct {
	Config                    *config.App
	AuthorizationDependencies *authorizationRouter.Dependencies
}

func initializeDependencies(path string, env string) (*Dependencies, error) {
	if err := loadEnvFile(env); err != nil {
		return nil, fmt.Errorf("load env failed,err:%w", err)
	}

	appConfig, err := loadConfigFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config file failed,err:%w", err)
	}

	if err := zap.Initialize(appConfig.Logger); err != nil {
		return nil, fmt.Errorf("initialize zap logger failed, err:%w", err)
	}

	mysqlDatabase, err := gorm.ConnectToMysql(appConfig.Mysql)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql failed, err:%w", err)
	}

	if err := gorm.AutoMigrate(mysqlDatabase, &authorizationModels.User{}, &models.Event{}); err != nil {
		return nil, fmt.Errorf("mysql mirgrate failed,err:%w", err)
	}

	redisClient, err := redis.ConnectToRedis(appConfig.Redis)
	if err != nil {
		return nil, fmt.Errorf("connect to redis failed, err:%w", err)
	}

	validator, err := ginutils.NewValidator()
	if err != nil {
		return nil, fmt.Errorf("validator initialize failed, err:%w", err)
	}

	numberGenerator, err := snowflake.NewNumberGenerator(appConfig.MachineNode)
	if err != nil {
		return nil, fmt.Errorf("number generator initialize failed, err:%w", err)
	}

	idGenerator := uuid.NewIDGenerator()

	hasher := bcrypt.NewHasher(appConfig.Hasher)

	authorizationUserRepository := authorizationRepository.NewUserRepository(mysqlDatabase, &authorizationConverters.UserConverter{})
	authorizationRefreshTokenRepository := authorizationRepository.NewRefreshTokenRepository(redisClient, &authorizationConverters.RefreshTokenConverter{})

	eventSaver := repository.NewEventSaver(mysqlDatabase, &converters.EventToModelConverter{})

	accessTokenManager := jwt.NewAccessTokenManager(appConfig.AccessToken)
	refreshTokenGenerator := crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)

	signUpUseCase := authorizationUseCase.NewSignUp(idGenerator, numberGenerator, hasher, authorizationUserRepository, eventSaver)
	loginUseCase := authorizationUseCase.NewLogin(idGenerator, hasher, authorizationUserRepository, authorizationRefreshTokenRepository, accessTokenManager, refreshTokenGenerator, eventSaver)
	refreshAccessTokenUseCase := authorizationUseCase.NewRefreshAccessToken(authorizationRefreshTokenRepository, authorizationRefreshTokenRepository, accessTokenManager, refreshTokenGenerator)
	parseAccessTokenUseCase := authorizationUseCase.NewParseAccessToken(accessTokenManager)

	authorizationMiddleware := middlewares.Authorization(parseAccessTokenUseCase)

	authorizationDependencies := &authorizationRouter.Dependencies{
		AuthorizationMiddleware:   authorizationMiddleware,
		SignUpUseCase:             signUpUseCase,
		LoginUseCase:              loginUseCase,
		RefreshAccessTokenUseCase: refreshAccessTokenUseCase,
		Validator:                 validator,
		RedisClient:               redisClient,
		RateLimitConfig:           appConfig.RateLimit,
		CookieConfig:              appConfig.Cookie,
	}

	return &Dependencies{
		Config:                    appConfig,
		AuthorizationDependencies: authorizationDependencies,
	}, nil
}
