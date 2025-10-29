package main

import (
	"fmt"
	"gochat/config"
	authorizationUsecase "gochat/internal/authorization/application/usecase"
	"gochat/internal/authorization/infrastructure/bcrypt"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authorizationConverters "gochat/internal/authorization/infrastructure/persistence/converter"
	authorizationModels "gochat/internal/authorization/infrastructure/persistence/model"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	"gochat/internal/authorization/infrastructure/snowflake"
	authorizationUuid "gochat/internal/authorization/infrastructure/uuid"
	"gochat/internal/infrastructure/canal"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/godotenv"
	gormutils "gochat/internal/infrastructure/gorm"
	kafkapub "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/persistence/converter"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/infrastructure/persistence/repository"
	redisutils "gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/uuid"
	"gochat/internal/infrastructure/viper"
	zaputils "gochat/internal/infrastructure/zap"

	"context"

	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	config                    *config.App
	signUpUseCase             authorizationUsecase.SignUpUseCase
	loginUseCase              authorizationUsecase.LoginUseCase
	refreshAccessTokenUseCase authorizationUsecase.RefreshAccessTokenUseCase
	parseAccessTokenUseCase   authorizationUsecase.ParseAccessTokenUseCase
	validator                 *ginutils.Validator
	redisClient               *redis.Client
	cancelAll                 context.CancelFunc
}

func initializeDependencies(path string, env string) (*Dependencies, error) {
	if err := godotenv.LoadEnvFile(env); err != nil {
		return nil, fmt.Errorf("load env failed,err:%w", err)
	}

	appConfig, err := viper.LoadConfigFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config file failed,err:%w", err)
	}

	if err := zaputils.Initialize(appConfig.Logger); err != nil {
		return nil, fmt.Errorf("initialize zap logger failed, err:%w", err)
	}

	mysqlDatabase, err := gormutils.ConnectToMysql(appConfig.Mysql)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql failed, err:%w", err)
	}

	if err := gormutils.AutoMigrate(mysqlDatabase, &authorizationModels.User{}, &model.Event{}, &model.DeadEvent{}); err != nil {
		return nil, fmt.Errorf("mysql mirgrate failed,err:%w", err)
	}

	redisClient, err := redisutils.ConnectToRedis(appConfig.Redis)
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

	eventIDGenerator := uuid.NewEventIDGenerator()
	userIDGenerator := authorizationUuid.NewUserIDGenerator()

	hasher := bcrypt.NewHasher(appConfig.Hasher)

	eventRepository := repository.NewEventRepository(mysqlDatabase, &converter.StandardEventConverter{})

	authorizationUserRepository := authorizationRepository.NewUserRepository(mysqlDatabase, &authorizationConverters.UserConverter{})
	authorizationRefreshTokenRepository := authorizationRepository.NewRefreshTokenRepository(redisClient, &authorizationConverters.RefreshTokenConverter{})

	accessTokenManager := jwt.NewAccessTokenManager(appConfig.AccessToken)
	refreshTokenGenerator := crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)

	signUpUseCase := authorizationUsecase.NewSignUpUseCase(eventIDGenerator, userIDGenerator, numberGenerator, hasher, authorizationUserRepository, eventRepository)
	loginUseCase := authorizationUsecase.NewLoginUseCase(eventIDGenerator, hasher, authorizationUserRepository, authorizationRefreshTokenRepository, accessTokenManager, refreshTokenGenerator, eventRepository)
	refreshAccessTokenUseCase := authorizationUsecase.NewRefreshAccessTokenUseCase(authorizationRefreshTokenRepository, authorizationRefreshTokenRepository, accessTokenManager, refreshTokenGenerator)
	parseAccessTokenUseCase := authorizationUsecase.NewParseAccessTokenUseCase(accessTokenManager)

	// start binlog outbox consumer
	ctx, cancel := context.WithCancel(context.Background())

	kafkaPublisher, err := kafkapub.NewEventPublisher(appConfig.Kafka, eventRepository, eventRepository)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("initialize kafka publisher failed, err:%w", err)
	}
	kafkaPublisher.Start()

	consumer, err := canal.NewOutboxConsumer(appConfig.BinlogReader, kafkaPublisher, eventRepository)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("initialize outbox consumer failed, err:%w", err)
	}

	consumer.Start(ctx)

	return &Dependencies{
		config:                    appConfig,
		signUpUseCase:             signUpUseCase,
		loginUseCase:              loginUseCase,
		refreshAccessTokenUseCase: refreshAccessTokenUseCase,
		parseAccessTokenUseCase:   parseAccessTokenUseCase,
		validator:                 validator,
		redisClient:               redisClient,
		cancelAll:                 cancel,
	}, nil
}
