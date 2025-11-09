package main

import (
	"context"
	"fmt"
	"gochat/config"
	authorizationUsecase "gochat/internal/authorization/application/usecase"
	authorizationDomain "gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authorizationConverter "gochat/internal/authorization/infrastructure/persistence/converter"
	authorizationModel "gochat/internal/authorization/infrastructure/persistence/model"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	authorizationSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authorizationUuid "gochat/internal/authorization/infrastructure/uuid"
	authorizationKafka "gochat/internal/authorization/port/kafka"
	chatUsecase "gochat/internal/chat/application/usecase"
	chatDomain "gochat/internal/chat/domain"
	chatModel "gochat/internal/chat/infrastructure/persistence/model"
	chatRepository "gochat/internal/chat/infrastructure/persistence/repository"
	chatUuid "gochat/internal/chat/infrastructure/uuid"
	chatKafka "gochat/internal/chat/port/kafka"
	"gochat/internal/infrastructure/bcrypt"
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
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/gomail"
	notificationConverter "gochat/internal/notification/infrastructure/persistence/converter"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	notificationUuid "gochat/internal/notification/infrastructure/uuid"
	notificationKafka "gochat/internal/notification/port/kafka"
	socialUseCase "gochat/internal/social/application/usecase"
	socialDomain "gochat/internal/social/domain"
	socialModel "gochat/internal/social/infrastructure/persistence/model"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"
	socialSnowflake "gochat/internal/social/infrastructure/snowflake"
	socialUuid "gochat/internal/social/infrastructure/uuid"
	socialKafka "gochat/internal/social/port/kafka"

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
	createRoomUseCase         socialUseCase.CreateRoomUseCase
	joinRoomUseCase           socialUseCase.JoinRoomUseCase
	leaveRoomUseCase          socialUseCase.LeaveRoomUseCase
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

	if appConfig.NeedMigrate {
		if err := gormutils.AutoMigrate(
			mysqlDatabase,
			&authorizationModel.User{},
			&notificationModel.Mail{},
			&notificationModel.User{},
			&notificationModel.Room{},
			&chatModel.User{},
			&chatModel.Room{},
			&chatModel.Message{},
			&chatModel.UserMessageState{},
			&socialModel.User{},
			&socialModel.Room{},
			&model.Event{},
			&model.DeadLetter{},
		); err != nil {
			return nil, fmt.Errorf("mysql mirgrate failed,err:%w", err)
		}
	}

	redisClient, err := redisutils.ConnectToRedis(appConfig.Redis)
	if err != nil {
		return nil, fmt.Errorf("connect to redis failed, err:%w", err)
	}

	validator, err := ginutils.NewValidator()
	if err != nil {
		return nil, fmt.Errorf("validator initialize failed, err:%w", err)
	}

	numberGenerator, err := authorizationSnowflake.NewUserNumberGenerator(appConfig.MachineNode)
	if err != nil {
		return nil, fmt.Errorf("number generator initialize failed, err:%w", err)
	}

	eventIDGenerator := uuid.NewEventIDGenerator()
	userIDGenerator := authorizationUuid.NewUserIDGenerator()

	hasher := bcrypt.NewHasher(appConfig.Hasher)

	eventRepository := repository.NewEventRepository(mysqlDatabase, &converter.StandardEventConverter{})

	authorizationUserRepository := authorizationRepository.NewUserRepository(mysqlDatabase, &authorizationConverter.UserConverter{})
	authorizationRefreshTokenRepository := authorizationRepository.NewRefreshTokenRepository(redisClient, &authorizationConverter.RefreshTokenConverter{})

	chatUserRepository := chatRepository.NewUserRepository(mysqlDatabase)
	chatRoomRepository := chatRepository.NewRoomRepository(mysqlDatabase)
	chatMessageRepository := chatRepository.NewMessageRepository(mysqlDatabase)

	socialUserRepository := socialRepository.NewUserRepository(mysqlDatabase)
	socialRoomRepository := socialRepository.NewRoomRepository(mysqlDatabase)

	notificationMailRepository := notificationRepository.NewMailRepository(mysqlDatabase, &notificationConverter.MailConverter{})
	notificationUserRepository := notificationRepository.NewUserRepository(mysqlDatabase)
	notificationRoomRepository := notificationRepository.NewRoomRepository(mysqlDatabase)

	accessTokenManager := jwt.NewAccessTokenManager(appConfig.AccessToken)
	refreshTokenGenerator := crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)

	messageNotifier := websocket.NewManager()

	emailNotifier := gomail.NewEmailNotifier(appConfig.Name, appConfig.Email)

	messageIDGenerator := chatUuid.NewMessageIDGenerator()

	mailIDGenerator := notificationUuid.NewMailIDGenerator()

	roomIDGenerator := socialUuid.NewRoomIDGenerator()

	roomNumberGenerator, err := socialSnowflake.NewRoomNumberGenerator(appConfig.MachineNode)
	if err != nil {
		return nil, err
	}

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
		authorizationUserRepository,
		eventRepository,
	)
	loginUseCase := authorizationUsecase.NewLoginUseCase(
		eventIDGenerator,
		hasher,
		authorizationUserRepository,
		authorizationUserRepository,
		authorizationRefreshTokenRepository,
		accessTokenManager,
		refreshTokenGenerator,
		eventRepository,
	)
	refreshAccessTokenUseCase := authorizationUsecase.NewRefreshAccessTokenUseCase(
		authorizationRefreshTokenRepository,
		authorizationRefreshTokenRepository,
		accessTokenManager,
		refreshTokenGenerator,
	)
	parseAccessTokenUseCase := authorizationUsecase.NewParseAccessTokenUseCase(accessTokenManager)

	sendPrivateMessageUseCase := chatUsecase.NewSendPrivateMessageUseCase(
		messageIDGenerator,
		eventIDGenerator,
		chatUserRepository,
		chatMessageRepository,
		eventRepository,
	)
	sendRoomMessageUseCase := chatUsecase.NewSendRoomMessageUseCase(
		messageIDGenerator,
		eventIDGenerator,
		chatRoomRepository,
		chatMessageRepository,
		eventRepository,
	)
	_ = chatUsecase.NewUpdateMessageStateUseCase(chatMessageRepository)

	_ = notificationUsecase.NewSendEmailUseCase(
		emailNotifier,
		notificationMailRepository,
		mailIDGenerator,
	)

	_ = notificationUsecase.NewSendMessageUseCase(
		eventIDGenerator,
		messageNotifier,
	)

	createRoomUseCase := socialUseCase.NewCreateRoomUseCase(
		eventIDGenerator,
		roomIDGenerator,
		roomNumberGenerator,
		hasher,
		socialRoomRepository,
		eventRepository,
	)
	joinRoomUseCase := socialUseCase.NewJoinRoomUseCase(
		eventIDGenerator,
		eventRepository,
		socialRoomRepository,
		hasher,
		socialRoomRepository,
	)
	leaveRoomUseCase := socialUseCase.NewLeaveRoomUseCase(
		eventIDGenerator,
		socialRoomRepository,
		socialRoomRepository,
		eventRepository,
	)

	// Initialize event-consumer usecases before subscriptions
	authUserCreatedUseCase := authorizationUsecase.NewUserCreatedUseCase(
		eventIDGenerator,
		kafkaPublisher,
		authorizationUserRepository,
	)

	socialUserCreatedUseCase := socialUseCase.NewUserCreatedUseCase(
		socialUserRepository,
	)
	socialRoomCreatedUseCase := socialUseCase.NewRoomCreatedUseCase(
		eventIDGenerator,
		kafkaPublisher,
		socialRoomRepository,
	)
	socialRoomJoinedUseCase := socialUseCase.NewRoomJoinedUseCase(
		eventIDGenerator,
		kafkaPublisher,
	)
	socialRoomLeftUseCase := socialUseCase.NewRoomLeftUseCase(
		eventIDGenerator,
		kafkaPublisher,
	)

	chatUserCreatedUseCase := chatUsecase.NewUserCreatedUseCase(
		chatUserRepository,
	)
	chatRoomCreatedUseCase := chatUsecase.NewRoomCreatedUseCase(
		chatRoomRepository,
	)
	chatRoomJoinedUseCase := chatUsecase.NewRoomJoinedUseCase(
		chatRoomRepository,
	)
	chatRoomLeftUseCase := chatUsecase.NewRoomLeftUseCase(
		chatRoomRepository,
	)

	notificationUserCreatedUseCase := notificationUsecase.NewUserCreatedUseCase(
		notificationUserRepository,
	)
	notificationRoomCreatedUseCase := notificationUsecase.NewRoomCreatedUseCase(
		notificationRoomRepository,
	)
	notificationRoomJoinedUseCase := notificationUsecase.NewRoomJoinedUseCase(
		notificationRoomRepository,
	)
	notificationRoomLeftUseCase := notificationUsecase.NewRoomLeftUseCase(
		notificationRoomRepository,
	)
	notificationPrivateMessageCreatedUseCase := notificationUsecase.NewPrivateMessageCreatedUseCase()
	notificationRoomMessageCreatedUseCase := notificationUsecase.NewRoomMessageCreatedUseCase()

	kafkaSubscriber, err := kafkautil.NewEventSubscriber(appConfig.Kafka.Common, appConfig.Kafka.Consumer)
	if err != nil {
		return nil, fmt.Errorf("initialize kafka kafkaSubscriber failed, err:%w", err)
	}

	// Authorization domain events -> usecase handler
	kafkaSubscriber.Subscribe(authorizationDomain.TopicUserCreated, authorizationKafka.NewUserCreatedEventHandler(
		authUserCreatedUseCase,
	))

	// Social domain events -> usecase handlers
	kafkaSubscriber.Subscribe(socialDomain.TopicUserCreated, socialKafka.NewUserCreatedEventHandler(
		socialUserCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(socialDomain.TopicRoomCreated, socialKafka.NewRoomCreatedEventHandler(
		socialRoomCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(socialDomain.TopicRoomJoined, socialKafka.NewRoomJoinedEventHandler(
		socialRoomJoinedUseCase,
	))
	kafkaSubscriber.Subscribe(socialDomain.TopicRoomLeft, socialKafka.NewRoomLeftEventHandler(
		socialRoomLeftUseCase,
	))

	// Chat domain events -> usecase handlers or bridging handlers
	kafkaSubscriber.Subscribe(chatDomain.TopicUserCreated, chatKafka.NewUserCreatedEventHandler(
		chatUserCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(chatDomain.TopicRoomCreated, chatKafka.NewRoomCreatedEventHandler(
		chatRoomCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(chatDomain.TopicRoomJoined, chatKafka.NewRoomJoinedEventHandler(
		chatRoomJoinedUseCase,
	))
	kafkaSubscriber.Subscribe(chatDomain.TopicRoomLeft, chatKafka.NewRoomLeftEventHandler(
		chatRoomLeftUseCase,
	))
	kafkaSubscriber.Subscribe(chatDomain.TopicPrivateMessageCreated, chatKafka.NewPrivateMessageCreatedEventHandler(
		eventIDGenerator,
		kafkaPublisher,
	))
	kafkaSubscriber.Subscribe(chatDomain.TopicRoomMessageCreated, chatKafka.NewRoomMessageCreatedEventHandler(
		eventIDGenerator,
		kafkaPublisher,
	))

	// Notification domain events -> usecase handlers
	kafkaSubscriber.Subscribe(notificationDomain.TopicUserCreated, notificationKafka.NewUserCreatedEventHandler(
		notificationUserCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(notificationDomain.TopicRoomCreated, notificationKafka.NewRoomCreatedEventHandler(
		notificationRoomCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(notificationDomain.TopicRoomJoined, notificationKafka.NewRoomJoinedEventHandler(
		notificationRoomJoinedUseCase,
	))
	kafkaSubscriber.Subscribe(notificationDomain.TopicRoomLeft, notificationKafka.NewRoomLeftEventHandler(
		notificationRoomLeftUseCase,
	))
	kafkaSubscriber.Subscribe(notificationDomain.TopicPrivateMessageCreated, notificationKafka.NewPrivateMessageCreatedEventHandler(
		notificationPrivateMessageCreatedUseCase,
	))
	kafkaSubscriber.Subscribe(notificationDomain.TopicRoomMessageCreated, notificationKafka.NewRoomMessageCreatedEventHandler(
		notificationRoomMessageCreatedUseCase,
	))

	// Start background components
	ctx, cancel := context.WithCancel(context.Background())

	if err := kafkaSubscriber.Start(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("start kafka kafkaSubscriber failed, err:%w", err)
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
		createRoomUseCase:         createRoomUseCase,
		joinRoomUseCase:           joinRoomUseCase,
		leaveRoomUseCase:          leaveRoomUseCase,
		validator:                 validator,
		redisClient:               redisClient,
	}, nil
}
