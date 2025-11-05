package main

import (
	"context"
	"fmt"
	"gochat/config"
	authorizationUsecase "gochat/internal/authorization/application/usecase"
	"gochat/internal/authorization/infrastructure/bcrypt"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authorizationConverter "gochat/internal/authorization/infrastructure/persistence/converter"
	authorizationModel "gochat/internal/authorization/infrastructure/persistence/model"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	"gochat/internal/authorization/infrastructure/snowflake"
	authorizationUuid "gochat/internal/authorization/infrastructure/uuid"
	chatUsecase "gochat/internal/chat/application/usecase"
	messageConverter "gochat/internal/chat/infrastructure/persistence/converter"
	chatModel "gochat/internal/chat/infrastructure/persistence/model"
	chatRepository "gochat/internal/chat/infrastructure/persistence/repository"
	chatUuid "gochat/internal/chat/infrastructure/uuid"
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
	"gochat/internal/infrastructure/websocket"
	zaputils "gochat/internal/infrastructure/zap"
	notificationUsecase "gochat/internal/notification/application/usecase"
	"gochat/internal/notification/infrastructure/gomail"
	notificationConverter "gochat/internal/notification/infrastructure/persistence/converter"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	notificationUuid "gochat/internal/notification/infrastructure/uuid"
	socialConverter "gochat/internal/social/infrastructure/persistence/converter"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"

	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	config                    *config.App
	closeAll                  func()
	signUpUseCase             authorizationUsecase.SignUpUseCase
	loginUseCase              authorizationUsecase.LoginUseCase
	refreshAccessTokenUseCase authorizationUsecase.RefreshAccessTokenUseCase
	parseAccessTokenUseCase   authorizationUsecase.ParseAccessTokenUseCase
	sendPrivateMessageUseCase chatUsecase.SendPrivateMessageUseCase
	sendRoomMessageUseCase    chatUsecase.SendRoomMessageUseCase
	validator                 *ginutils.Validator
	redisClient               *redis.Client
}

func initializeDependencies(configPath string, envPath string) (*Dependencies, error) {
	if err := godotenv.LoadEnvFile(envPath); err != nil {
		return nil, fmt.Errorf("load env failed,err:%w", err)
	}

	appConfig, err := viper.LoadConfigFile(configPath)
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
		&authorizationModel.User{},
		&notificationModel.Mail{},
		&chatModel.Message{},
		&chatModel.RecipientMessageState{},
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

	userRepository := authorizationRepository.NewUserRepository(mysqlDatabase, &authorizationConverter.UserConverter{})
	refreshTokenRepository := authorizationRepository.NewRefreshTokenRepository(redisClient, &authorizationConverter.RefreshTokenConverter{})

	messageRepository := chatRepository.NewMessageRepository(mysqlDatabase, &messageConverter.MessageConverter{})

	roomRepository := socialRepository.NewRoomRepository(mysqlDatabase, &socialConverter.RoomConverter{})

	mailRepository := notificationRepository.NewMailRepository(mysqlDatabase, &notificationConverter.MailConverter{})

	accessTokenManager := jwt.NewAccessTokenManager(appConfig.AccessToken)
	refreshTokenGenerator := crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)

	messageNotifier := websocket.NewManager()

	emailNotifier := gomail.NewEmailNotifier(appConfig.Name, appConfig.Email)

	messageIDGenerator := chatUuid.NewMessageIDGenerator()

	mailIDGenerator := notificationUuid.NewMailIDGenerator()

	kafkaPublisher, err := kafkautil.NewEventPublisher(appConfig.Kafka.Common, appConfig.Kafka.Producer, eventRepository, eventRepository)
	if err != nil {
		return nil, fmt.Errorf("initialize kafka publisher failed, err:%w", err)
	}

	consumer, err := canal.NewOutboxConsumer(appConfig.BinlogReader, kafkaPublisher, eventRepository)
	if err != nil {
		return nil, fmt.Errorf("initialize outbox consumer failed, err:%w", err)
	}

	signUpUseCase := authorizationUsecase.NewSignUpUseCase(
		eventIDGenerator,
		userIDGenerator,
		numberGenerator,
		hasher,
		userRepository,
		eventRepository,
	)
	loginUseCase := authorizationUsecase.NewLoginUseCase(
		eventIDGenerator,
		hasher,
		userRepository,
		refreshTokenRepository,
		accessTokenManager,
		refreshTokenGenerator,
		eventRepository,
	)
	refreshAccessTokenUseCase := authorizationUsecase.NewRefreshAccessTokenUseCase(
		refreshTokenRepository,
		refreshTokenRepository,
		accessTokenManager,
		refreshTokenGenerator,
	)
	parseAccessTokenUseCase := authorizationUsecase.NewParseAccessTokenUseCase(accessTokenManager)

	sendPrivateMessageUseCase := chatUsecase.NewSendPrivateMessageUseCase(
		messageIDGenerator,
		eventIDGenerator,
		userRepository,
		messageRepository,
		eventRepository,
	)
	sendRoomMessageUseCase := chatUsecase.NewSendRoomMessageUseCase(
		messageIDGenerator,
		eventIDGenerator,
		roomRepository,
		messageRepository,
		eventRepository,
	)
	updateMessageStateUseCase := chatUsecase.NewUpdateMessageStateUseCase(messageRepository)

	sendEmailUseCase := notificationUsecase.NewSendEmailUseCase(emailNotifier, mailRepository, mailIDGenerator)

	sendMessageUseCase := notificationUsecase.NewSendMessageUseCase(eventIDGenerator, messageNotifier, kafkaPublisher)

	kafkaSubscriber, err := kafka.NewSubscriber(
		appConfig.Kafka.Common,
		appConfig.Kafka.Consumer,
		eventIDGenerator,
		kafkaPublisher,
		updateMessageStateUseCase,
		sendEmailUseCase,
		sendMessageUseCase,
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
		sendPrivateMessageUseCase: sendPrivateMessageUseCase,
		sendRoomMessageUseCase:    sendRoomMessageUseCase,
		validator:                 validator,
		redisClient:               redisClient,
	}, nil
}
