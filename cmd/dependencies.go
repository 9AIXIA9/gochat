package main

import (
	"context"
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
	"gochat/internal/delivery/kafka"
	"gochat/internal/infrastructure/canal"
	"gochat/internal/infrastructure/godotenv"
	gormutils "gochat/internal/infrastructure/gorm"
	kafkautil "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/persistence/converter"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/infrastructure/persistence/repository"
	redisutils "gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/uuid"
	ginutils "gochat/internal/infrastructure/validator"
	"gochat/internal/infrastructure/viper"
	zaputils "gochat/internal/infrastructure/zap"
	notificationUsecase "gochat/internal/notification/application/usecase"
	"gochat/internal/notification/infrastructure/gomail"
	notificationConverters "gochat/internal/notification/infrastructure/persistence/converter"
	notificationModels "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	notificationUuid "gochat/internal/notification/infrastructure/uuid"

	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	config                    *config.App
	closeAll                  func()
	signUpUseCase             authorizationUsecase.SignUpUseCase
	loginUseCase              authorizationUsecase.LoginUseCase
	refreshAccessTokenUseCase authorizationUsecase.RefreshAccessTokenUseCase
	parseAccessTokenUseCase   authorizationUsecase.ParseAccessTokenUseCase
	validator                 *ginutils.Validator
	redisClient               *redis.Client
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

	if err := gormutils.AutoMigrate(
		mysqlDatabase,
		&authorizationModels.User{},
		&notificationModels.Notice{},
		&model.Event{},
		&model.DeadLetter{},
	); err != nil {
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

	noticeRepository := notificationRepository.NewNoticeRepository(mysqlDatabase, &notificationConverters.NoticeConverter{})

	accessTokenManager := jwt.NewAccessTokenManager(appConfig.AccessToken)
	refreshTokenGenerator := crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)

	emailNotifier := gomail.NewEmailNotifier(appConfig.Name, appConfig.Email)

	noticeIDGenerator := notificationUuid.NewNoticeIDGenerator()

	signUpUseCase := authorizationUsecase.NewSignUpUseCase(eventIDGenerator, userIDGenerator, numberGenerator, hasher, authorizationUserRepository, eventRepository)
	loginUseCase := authorizationUsecase.NewLoginUseCase(eventIDGenerator, hasher, authorizationUserRepository, authorizationRefreshTokenRepository, accessTokenManager, refreshTokenGenerator, eventRepository)
	refreshAccessTokenUseCase := authorizationUsecase.NewRefreshAccessTokenUseCase(authorizationRefreshTokenRepository, authorizationRefreshTokenRepository, accessTokenManager, refreshTokenGenerator)
	parseAccessTokenUseCase := authorizationUsecase.NewParseAccessTokenUseCase(accessTokenManager)

	sendEmailUseCase := notificationUsecase.NewSendEmailUseCase(emailNotifier, noticeRepository, noticeIDGenerator)

	kafkaPublisher, err := kafkautil.NewEventPublisher(appConfig.Kafka.Common, appConfig.Kafka.Producer, eventRepository, eventRepository)
	if err != nil {
		return nil, fmt.Errorf("initialize kafka publisher failed, err:%w", err)
	}

	consumer, err := canal.NewOutboxConsumer(appConfig.BinlogReader, kafkaPublisher, eventRepository)
	if err != nil {
		return nil, fmt.Errorf("initialize outbox consumer failed, err:%w", err)
	}

	kafkaSubscriber, err := kafka.NewSubscriber(
		appConfig.Kafka.Common,
		appConfig.Kafka.Consumer,
		kafkaPublisher,
		sendEmailUseCase,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize kafka subscriber failed, err:%w", err)
	}

	// Start background components
	ctx, cancel := context.WithCancel(context.Background())

	if err := kafkaSubscriber.Start(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("start kafka subscriber failed, err:%w", err)
	}
	consumer.Start()
	kafkaPublisher.Start()
	emailNotifier.Start()

	return &Dependencies{
		config: appConfig,
		closeAll: func() {
			consumer.Close()
			kafkaPublisher.Close()
			emailNotifier.Close()
			kafkaSubscriber.Close()
			cancel()
		},
		signUpUseCase:             signUpUseCase,
		loginUseCase:              loginUseCase,
		refreshAccessTokenUseCase: refreshAccessTokenUseCase,
		parseAccessTokenUseCase:   parseAccessTokenUseCase,
		validator:                 validator,
		redisClient:               redisClient,
	}, nil
}
