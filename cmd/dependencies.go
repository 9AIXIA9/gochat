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
	zaputils "gochat/internal/infrastructure/zap"
	notificationUsecase "gochat/internal/notification/application/usecase"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/gomail"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	notificationKafka "gochat/internal/notification/port/kafka"
	"gochat/internal/shared/event"
	socialUseCase "gochat/internal/social/application/usecase"
	socialDomain "gochat/internal/social/domain"
	socialModel "gochat/internal/social/infrastructure/persistence/model"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"
	socialSnowflake "gochat/internal/social/infrastructure/snowflake"
	socialUuid "gochat/internal/social/infrastructure/uuid"
	socialKafka "gochat/internal/social/port/kafka"

	authorizationApp "gochat/internal/authorization/application"
	chatApp "gochat/internal/chat/application"
	socialApp "gochat/internal/social/application"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConfigPath string

type EnvPath string

// Dependencies holds runtime objects & usecases.
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

// -------------------- Base Providers --------------------
func provideAppConfig(configPath ConfigPath, envPath EnvPath) (*config.App, error) {
	if err := godotenv.LoadEnvFile(string(envPath)); err != nil {
		return nil, fmt.Errorf("load env failed,err:%w", err)
	}
	cfg, err := viper.LoadConfigFile(string(configPath))
	if err != nil {
		return nil, fmt.Errorf("load config file failed,err:%w", err)
	}
	if err := zaputils.Initialize(cfg.Logger); err != nil {
		return nil, fmt.Errorf("initialize zap logger failed, err:%w", err)
	}
	return cfg, nil
}
func provideMysql(appConfig *config.App) (*gorm.DB, error) {
	return gormutils.ConnectToMysql(appConfig.Mysql)
}
func provideRedis(appConfig *config.App) (*redis.Client, error) {
	return redisutils.ConnectToRedis(appConfig.Redis)
}
func provideValidator() (*ginutils.Validator, error) { return ginutils.NewValidator() }

// -------------------- Generators & Managers --------------------
func provideEventIDGenerator() *uuid.EventIDGenerator { return uuid.NewEventIDGenerator() }
func provideAuthorizationUserIDGenerator() *authorizationUuid.UserIDGenerator {
	return authorizationUuid.NewUserIDGenerator()
}
func provideAuthorizationUserNumberGenerator(appConfig *config.App) (*authorizationSnowflake.UserNumberGenerator, error) {
	return authorizationSnowflake.NewUserNumberGenerator(appConfig.MachineNode)
}
func provideHasher(appConfig *config.App) *bcrypt.Hasher { return bcrypt.NewHasher(appConfig.Hasher) }
func provideMessageIDGenerator() *chatUuid.MessageIDGenerator {
	return chatUuid.NewMessageIDGenerator()
}

func provideSocialRoomIDGenerator() *socialUuid.RoomIDGenerator {
	return socialUuid.NewRoomIDGenerator()
}
func provideSocialRoomNumberGenerator(appConfig *config.App) (*socialSnowflake.RoomNumberGenerator, error) {
	return socialSnowflake.NewRoomNumberGenerator(appConfig.MachineNode)
}
func provideAccessTokenManager(appConfig *config.App) *jwt.AccessTokenManager {
	return jwt.NewAccessTokenManager(appConfig.AccessToken)
}
func provideRefreshTokenGenerator(appConfig *config.App) *crypto.RefreshTokenGenerator {
	return crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)
}
func provideEmailNotifier(appConfig *config.App) *gomail.EmailNotifier {
	return gomail.NewEmailNotifier(appConfig.Name, appConfig.Email)
}

// -------------------- Repositories --------------------
func provideEventRepository(mysql *gorm.DB) *repository.EventRepository {
	return repository.NewEventRepository(mysql, &converter.StandardEventConverter{})
}
func provideAuthorizationUserRepository(mysql *gorm.DB) *authorizationRepository.UserRepository {
	return authorizationRepository.NewUserRepository(mysql, &authorizationConverter.UserConverter{})
}
func provideAuthorizationRefreshTokenRepository(redisClient *redis.Client) *authorizationRepository.RefreshTokenRepository {
	return authorizationRepository.NewRefreshTokenRepository(redisClient, &authorizationConverter.RefreshTokenConverter{})
}
func provideChatUserRepository(mysql *gorm.DB) *chatRepository.UserRepository {
	return chatRepository.NewUserRepository(mysql)
}
func provideChatRoomRepository(mysql *gorm.DB) *chatRepository.RoomRepository {
	return chatRepository.NewRoomRepository(mysql)
}
func provideChatMessageRepository(mysql *gorm.DB) *chatRepository.MessageRepository {
	return chatRepository.NewMessageRepository(mysql)
}
func provideSocialUserRepository(mysql *gorm.DB) *socialRepository.UserRepository {
	return socialRepository.NewUserRepository(mysql)
}
func provideSocialRoomRepository(mysql *gorm.DB) *socialRepository.RoomRepository {
	return socialRepository.NewRoomRepository(mysql)
}
func provideNotificationUserRepository(mysql *gorm.DB) *notificationRepository.UserRepository {
	return notificationRepository.NewUserRepository(mysql)
}
func provideNotificationRoomRepository(mysql *gorm.DB) *notificationRepository.RoomRepository {
	return notificationRepository.NewRoomRepository(mysql)
}

// -------------------- Kafka & Canal --------------------
func provideKafkaPublisher(appConfig *config.App, eventRepo *repository.EventRepository) (*kafkautil.EventPublisher, error) {
	return kafkautil.NewEventPublisher(appConfig.Kafka.Common, appConfig.Kafka.Producer, eventRepo, eventRepo)
}
func provideKafkaSubscriber(appConfig *config.App) (*kafkautil.EventSubscriber, error) {
	return kafkautil.NewEventSubscriber(appConfig.Kafka.Common, appConfig.Kafka.Consumer)
}
func provideOutboxConsumer(appConfig *config.App, publisher *kafkautil.EventPublisher, eventRepo *repository.EventRepository) (*canal.OutboxConsumer, error) {
	return canal.NewOutboxConsumer(appConfig.BinlogReader, publisher, eventRepo)
}

// -------------------- UseCases (HTTP side) --------------------
func provideSignUpUseCase(
	eventIDGen event.IDGenerator,
	userIDGen authorizationApp.UserIDGenerator,
	numberGen authorizationApp.UserNumberGenerator,
	encryptor authorizationApp.Encryptor,
	userSaver authorizationApp.UserSaver,
	eventSaver event.UnpublishedSaver,
) authorizationUsecase.SignUpUseCase {
	return authorizationUsecase.NewSignUpUseCase(eventIDGen, userIDGen, numberGen, encryptor, userSaver, eventSaver)
}

func provideLoginUseCase(
	eventIDGen event.IDGenerator,
	comparator authorizationApp.Comparator,
	userFinder authorizationApp.UserFinderByNumber,
	userUpdater authorizationApp.UserLoggedInAtUpdater,
	refreshTokenSaver authorizationApp.RefreshTokenSaver,
	accessTokenGenerator authorizationApp.AccessTokenGenerator,
	refreshTokenGenerator authorizationApp.RefreshTokenGenerator,
	eventSaver event.UnpublishedSaver,
) authorizationUsecase.LoginUseCase {
	return authorizationUsecase.NewLoginUseCase(eventIDGen, comparator, userFinder, userUpdater, refreshTokenSaver, accessTokenGenerator, refreshTokenGenerator, eventSaver)
}
func provideRefreshAccessTokenUseCase(
	refreshTokenSaver authorizationApp.RefreshTokenSaver,
	refreshTokenFinder authorizationApp.RefreshTokenFinder,
	accessTokenGenerator authorizationApp.AccessTokenGenerator,
	refreshTokenGenerator authorizationApp.RefreshTokenGenerator,
) authorizationUsecase.RefreshAccessTokenUseCase {
	return authorizationUsecase.NewRefreshAccessTokenUseCase(refreshTokenSaver, refreshTokenFinder, accessTokenGenerator, refreshTokenGenerator)
}
func provideParseAccessTokenUseCase(accessTokenParser authorizationApp.AccessTokenParser) authorizationUsecase.ParseAccessTokenUseCase {
	return authorizationUsecase.NewParseAccessTokenUseCase(accessTokenParser)
}
func provideSendPrivateMessageUseCase(
	messageIDGen chatApp.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	userRepo *chatRepository.UserRepository,
	msgRepo *chatRepository.MessageRepository,
	eventSaver event.UnpublishedSaver,
) chatUsecase.SendPrivateMessageUseCase {
	return chatUsecase.NewSendPrivateMessageUseCase(messageIDGen, eventIDGen, userRepo, msgRepo, eventSaver)
}
func provideSendRoomMessageUseCase(
	messageIDGen chatApp.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	roomRepo *chatRepository.RoomRepository,
	msgRepo *chatRepository.MessageRepository,
	eventSaver event.UnpublishedSaver,
) chatUsecase.SendRoomMessageUseCase {
	return chatUsecase.NewSendRoomMessageUseCase(messageIDGen, eventIDGen, roomRepo, msgRepo, eventSaver)
}
func provideCreateRoomUseCase(
	eventIDGen event.IDGenerator,
	roomIDGen socialApp.RoomIDGenerator,
	numberGen socialApp.RoomNumberGenerator,
	encryptor socialApp.Encryptor,
	roomSaver socialApp.RoomSaver,
	eventSaver event.UnpublishedSaver,
) socialUseCase.CreateRoomUseCase {
	return socialUseCase.NewCreateRoomUseCase(eventIDGen, roomIDGen, numberGen, encryptor, roomSaver, eventSaver)
}
func provideJoinRoomUseCase(
	eventIDGen event.IDGenerator,
	eventSaver event.UnpublishedSaver,
	finder socialApp.RoomFinder,
	comparator socialApp.Comparator,
	roomJoiner socialApp.RoomJoiner,
) socialUseCase.JoinRoomUseCase {
	return socialUseCase.NewJoinRoomUseCase(eventIDGen, eventSaver, finder, comparator, roomJoiner)
}
func provideLeaveRoomUseCase(
	eventIDGen event.IDGenerator,
	finder socialApp.RoomFinder,
	roomLeaver socialApp.RoomLeaver,
	eventSaver event.UnpublishedSaver,
) socialUseCase.LeaveRoomUseCase {
	return socialUseCase.NewLeaveRoomUseCase(eventIDGen, finder, roomLeaver, eventSaver)
}

// -------------------- Event UseCases (Kafka consumer side) --------------------
func provideAuthUserCreatedUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher, userRepo *authorizationRepository.UserRepository) authorizationUsecase.UserCreatedUseCase {
	return authorizationUsecase.NewUserCreatedUseCase(eventIDGen, publisher, userRepo)
}
func provideSocialUserCreatedUseCase(userRepo *socialRepository.UserRepository) socialUseCase.UserCreatedUseCase {
	return socialUseCase.NewUserCreatedUseCase(userRepo)
}
func provideSocialRoomCreatedUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher, roomRepo *socialRepository.RoomRepository) socialUseCase.RoomCreatedUseCase {
	return socialUseCase.NewRoomCreatedUseCase(eventIDGen, publisher, roomRepo)
}
func provideSocialRoomJoinedUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher) socialUseCase.RoomJoinedUseCase {
	return socialUseCase.NewRoomJoinedUseCase(eventIDGen, publisher)
}
func provideSocialRoomLeftUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher) socialUseCase.RoomLeftUseCase {
	return socialUseCase.NewRoomLeftUseCase(eventIDGen, publisher)
}
func provideChatUserCreatedUseCase(userRepo *chatRepository.UserRepository) chatUsecase.UserCreatedUseCase {
	return chatUsecase.NewUserCreatedUseCase(userRepo)
}
func provideChatRoomCreatedUseCase(roomRepo *chatRepository.RoomRepository) chatUsecase.RoomCreatedUseCase {
	return chatUsecase.NewRoomCreatedUseCase(roomRepo)
}
func provideChatRoomJoinedUseCase(roomRepo *chatRepository.RoomRepository) chatUsecase.RoomJoinedUseCase {
	return chatUsecase.NewRoomJoinedUseCase(roomRepo)
}
func provideChatRoomLeftUseCase(roomRepo *chatRepository.RoomRepository) chatUsecase.RoomLeftUseCase {
	return chatUsecase.NewRoomLeftUseCase(roomRepo)
}
func provideNotificationUserCreatedUseCase(userRepo *notificationRepository.UserRepository, emailNotifier *gomail.EmailNotifier) notificationUsecase.UserCreatedUseCase {
	return notificationUsecase.NewUserCreatedUseCase(userRepo, emailNotifier)
}
func provideNotificationRoomCreatedUseCase(roomRepo *notificationRepository.RoomRepository) notificationUsecase.RoomCreatedUseCase {
	return notificationUsecase.NewRoomCreatedUseCase(roomRepo)
}
func provideNotificationRoomJoinedUseCase(roomRepo *notificationRepository.RoomRepository) notificationUsecase.RoomJoinedUseCase {
	return notificationUsecase.NewRoomJoinedUseCase(roomRepo)
}
func provideNotificationRoomLeftUseCase(roomRepo *notificationRepository.RoomRepository) notificationUsecase.RoomLeftUseCase {
	return notificationUsecase.NewRoomLeftUseCase(roomRepo)
}
func provideNotificationPrivateMessageCreatedUseCase() notificationUsecase.PrivateMessageCreatedUseCase {
	return notificationUsecase.NewPrivateMessageCreatedUseCase()
}
func provideNotificationRoomMessageCreatedUseCase() notificationUsecase.RoomMessageCreatedUseCase {
	return notificationUsecase.NewRoomMessageCreatedUseCase()
}

// -------------------- Kafka Subscriptions --------------------
func provideKafkaSubscriptions(subscriber *kafkautil.EventSubscriber,
	// auth
	authUserCreated authorizationUsecase.UserCreatedUseCase,
	// social
	socialUserCreated socialUseCase.UserCreatedUseCase,
	socialRoomCreated socialUseCase.RoomCreatedUseCase,
	socialRoomJoined socialUseCase.RoomJoinedUseCase,
	socialRoomLeft socialUseCase.RoomLeftUseCase,
	// chat
	chatUserCreated chatUsecase.UserCreatedUseCase,
	chatRoomCreated chatUsecase.RoomCreatedUseCase,
	chatRoomJoined chatUsecase.RoomJoinedUseCase,
	chatRoomLeft chatUsecase.RoomLeftUseCase,
	// notification
	notificationUserCreated notificationUsecase.UserCreatedUseCase,
	notificationRoomCreated notificationUsecase.RoomCreatedUseCase,
	notificationRoomJoined notificationUsecase.RoomJoinedUseCase,
	notificationRoomLeft notificationUsecase.RoomLeftUseCase,
	notificationPrivateMessageCreated notificationUsecase.PrivateMessageCreatedUseCase,
	notificationRoomMessageCreated notificationUsecase.RoomMessageCreatedUseCase,
	// bridging handlers (chat events require eventID & publisher)
	eventIDGen *uuid.EventIDGenerator,
	publisher *kafkautil.EventPublisher,
) error {
	// Authorization
	subscriber.Subscribe(authorizationDomain.TopicUserCreated, authorizationKafka.NewUserCreatedEventHandler(authUserCreated))
	// Social
	subscriber.Subscribe(socialDomain.TopicUserCreated, socialKafka.NewUserCreatedEventHandler(socialUserCreated))
	subscriber.Subscribe(socialDomain.TopicRoomCreated, socialKafka.NewRoomCreatedEventHandler(socialRoomCreated))
	subscriber.Subscribe(socialDomain.TopicRoomJoined, socialKafka.NewRoomJoinedEventHandler(socialRoomJoined))
	subscriber.Subscribe(socialDomain.TopicRoomLeft, socialKafka.NewRoomLeftEventHandler(socialRoomLeft))
	// Chat
	subscriber.Subscribe(chatDomain.TopicUserCreated, chatKafka.NewUserCreatedEventHandler(chatUserCreated))
	subscriber.Subscribe(chatDomain.TopicRoomCreated, chatKafka.NewRoomCreatedEventHandler(chatRoomCreated))
	subscriber.Subscribe(chatDomain.TopicRoomJoined, chatKafka.NewRoomJoinedEventHandler(chatRoomJoined))
	subscriber.Subscribe(chatDomain.TopicRoomLeft, chatKafka.NewRoomLeftEventHandler(chatRoomLeft))
	subscriber.Subscribe(chatDomain.TopicPrivateMessageCreated, chatKafka.NewPrivateMessageCreatedEventHandler(eventIDGen, publisher))
	subscriber.Subscribe(chatDomain.TopicRoomMessageCreated, chatKafka.NewRoomMessageCreatedEventHandler(eventIDGen, publisher))
	// Notification
	subscriber.Subscribe(notificationDomain.TopicUserCreated, notificationKafka.NewUserCreatedEventHandler(notificationUserCreated))
	subscriber.Subscribe(notificationDomain.TopicRoomCreated, notificationKafka.NewRoomCreatedEventHandler(notificationRoomCreated))
	subscriber.Subscribe(notificationDomain.TopicRoomJoined, notificationKafka.NewRoomJoinedEventHandler(notificationRoomJoined))
	subscriber.Subscribe(notificationDomain.TopicRoomLeft, notificationKafka.NewRoomLeftEventHandler(notificationRoomLeft))
	subscriber.Subscribe(notificationDomain.TopicPrivateMessageCreated, notificationKafka.NewPrivateMessageCreatedEventHandler(notificationPrivateMessageCreated))
	subscriber.Subscribe(notificationDomain.TopicRoomMessageCreated, notificationKafka.NewRoomMessageCreatedEventHandler(notificationRoomMessageCreated))
	return nil
}

func BuildDependencies(
	appConfig *config.App,
	mysql *gorm.DB,
	redisClient *redis.Client,
	validator *ginutils.Validator,
	// core infra
	kafkaPublisher *kafkautil.EventPublisher,
	kafkaSubscriber *kafkautil.EventSubscriber,
	outboxConsumer *canal.OutboxConsumer,
	emailNotifier *gomail.EmailNotifier,
	// usecases
	signUp authorizationUsecase.SignUpUseCase,
	login authorizationUsecase.LoginUseCase,
	refreshAccess authorizationUsecase.RefreshAccessTokenUseCase,
	parseAccess authorizationUsecase.ParseAccessTokenUseCase,
	sendPrivate chatUsecase.SendPrivateMessageUseCase,
	sendRoom chatUsecase.SendRoomMessageUseCase,
	createRoom socialUseCase.CreateRoomUseCase,
	joinRoom socialUseCase.JoinRoomUseCase,
	leaveRoom socialUseCase.LeaveRoomUseCase,
	// subscriptions side-effect
	_ error, // ensure subscriptions provider executed (ignored)
) (*Dependencies, error) {
	// Migrations (side-effect). Performed here to keep initialize logic centralized.
	if appConfig.NeedMigrate {
		if err := gormutils.AutoMigrate(
			mysql,
			&authorizationModel.User{},
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
			return nil, fmt.Errorf("mysql migrate failed, err:%w", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := kafkaSubscriber.Start(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("start kafka subscriber failed, err:%w", err)
	}
	outboxConsumer.Start()
	kafkaPublisher.Start()
	emailNotifier.Start()

	deps := &Dependencies{
		config: appConfig,
		closeAll: func() {
			outboxConsumer.Close()
			kafkaPublisher.Close()
			emailNotifier.Close()
			kafkaSubscriber.Close()
			cancel()
		},
		signUpUseCase:             signUp,
		loginUseCase:              login,
		refreshAccessTokenUseCase: refreshAccess,
		parseAccessTokenUseCase:   parseAccess,
		sendPrivateMessageUseCase: sendPrivate,
		sendRoomMessageUseCase:    sendRoom,
		createRoomUseCase:         createRoom,
		joinRoomUseCase:           joinRoom,
		leaveRoomUseCase:          leaveRoom,
		validator:                 validator,
		redisClient:               redisClient,
	}
	return deps, nil
}
