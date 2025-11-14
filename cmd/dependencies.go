package main

import (
	"context"
	"fmt"
	"gochat/config"
	"gochat/internal/application/usecase"
	authorizationUsecase "gochat/internal/authorization/application/usecase"
	authorizationDomain "gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authorizationConverter "gochat/internal/authorization/infrastructure/persistence/converter"
	authorizationModel "gochat/internal/authorization/infrastructure/persistence/model"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	authorizationSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authorizationUuid "gochat/internal/authorization/infrastructure/uuid"
	authorizationHttp "gochat/internal/authorization/port/http"
	authorizationKafka "gochat/internal/authorization/port/kafka"
	chatUsecase "gochat/internal/chat/application/usecase"
	chatDomain "gochat/internal/chat/domain"
	chatModel "gochat/internal/chat/infrastructure/persistence/model"
	chatRepository "gochat/internal/chat/infrastructure/persistence/repository"
	chatUuid "gochat/internal/chat/infrastructure/uuid"
	chatHttp "gochat/internal/chat/port/http"
	chatKafka "gochat/internal/chat/port/kafka"
	"gochat/internal/delivery/http/handler"
	"gochat/internal/delivery/http/middleware"
	"gochat/internal/delivery/kafka"
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
	gomailUtil "gochat/internal/notification/infrastructure/gomail"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	notificationWebsocketInfrastructure "gochat/internal/notification/infrastructure/websocket"
	notificationKafka "gochat/internal/notification/port/kafka"
	notificationWebsocket "gochat/internal/notification/port/websocket"
	"gochat/internal/shared/event"
	socialUseCase "gochat/internal/social/application/usecase"
	socialDomain "gochat/internal/social/domain"
	socialModel "gochat/internal/social/infrastructure/persistence/model"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"
	socialSnowflake "gochat/internal/social/infrastructure/snowflake"
	socialUuid "gochat/internal/social/infrastructure/uuid"
	socialHttp "gochat/internal/social/port/http"
	socialKafka "gochat/internal/social/port/kafka"
	"time"

	authorizationApp "gochat/internal/authorization/application"
	chatApp "gochat/internal/chat/application"
	socialApp "gochat/internal/social/application"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type ConfigPath string

type EnvPath string

// Dependencies holds runtime objects & usecases.
type Dependencies struct {
	HttpRouter *gin.Engine
	config     *config.App
	closeAll   func()
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

// -------------------- Generators & Notifiers --------------------
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
func provideGomailDialer(appConfig *config.App) (*gomail.Dialer, error) {
	return gomailUtil.NewDialer(appConfig.Email)
}
func provideEmailNotifier(appConfig *config.App, dialer *gomail.Dialer) *gomailUtil.EmailNotifier {
	return gomailUtil.NewEmailNotifier(appConfig.Name, dialer)
}
func provideMessageNotifier(manager *websocket.Manager) *notificationWebsocketInfrastructure.MessageNotifier {
	return notificationWebsocketInfrastructure.NewMessageNotifier(manager)
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
func provideNotificationMessageRepository(mysql *gorm.DB) *notificationRepository.MessageRepository {
	return notificationRepository.NewMessageRepository(mysql)
}

// -------------------- Kafka & Canal & Websocket --------------------
func provideKafkaConsumer(appConfig *config.App) (*ckafka.Consumer, error) {
	return kafkautil.NewConsumer(appConfig.Kafka.Common, appConfig.Kafka.Consumer)
}
func provideKafkaProducer(appConfig *config.App) (*ckafka.Producer, error) {
	return kafkautil.NewProducer(appConfig.Kafka.Common, appConfig.Kafka.Producer)
}
func provideKafkaProducerWithRetry(producer *ckafka.Producer) *kafkautil.ProducerWithRetry {
	return kafkautil.NewProducerWithRetry(producer)
}
func provideKafkaRetrier(producer *kafkautil.ProducerWithRetry, eventRepo *repository.EventRepository) *kafkautil.ConsumerRetrier {
	return kafkautil.NewConsumerRetrier(producer, eventRepo)
}
func provideKafkaPublisher(producer *kafkautil.ProducerWithRetry, eventRepo *repository.EventRepository) (*kafkautil.EventPublisher, error) {
	return kafkautil.NewEventPublisher(producer, eventRepo)
}
func provideKafkaSubscriber(consumer *ckafka.Consumer, retrier *kafkautil.ConsumerRetrier) (*kafkautil.EventSubscriber, error) {
	return kafkautil.NewEventSubscriber(consumer, retrier)
}

func provideOutboxConsumer(appConfig *config.App, publisher *kafkautil.EventPublisher, eventRepo *repository.EventRepository) (*canal.OutboxConsumer, error) {
	return canal.NewOutboxConsumer(appConfig.BinlogReader, publisher, eventRepo)
}
func provideWebsocketUpgrader(appConfig *config.App) *gorillaWebsocket.Upgrader {
	return websocket.NewUpgrader(appConfig.CORS.AllowOrigins)
}
func provideWebsocketManager(upgrader *gorillaWebsocket.Upgrader) *websocket.Manager {
	return websocket.NewManager(upgrader)
}
func provideWebsocketServer(manager *websocket.Manager, router *websocket.Router, publisher *kafkautil.EventPublisher, eventIDGen event.IDGenerator) *websocket.Server {
	return websocket.NewServer(manager, router, publisher, eventIDGen)
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

func provideKafkaTopics() []string {
	return []string{
		// websocket
		string(websocket.TopicUserSessionStarted),
		// authorization
		string(authorizationDomain.TopicUserCreated),
		// social
		string(socialDomain.TopicUserCreated),
		string(socialDomain.TopicRoomCreated),
		string(socialDomain.TopicRoomJoined),
		string(socialDomain.TopicRoomLeft),
		// chat
		string(chatDomain.TopicUserCreated),
		string(chatDomain.TopicRoomCreated),
		string(chatDomain.TopicRoomJoined),
		string(chatDomain.TopicRoomLeft),
		string(chatDomain.TopicPrivateMessageCreated),
		string(chatDomain.TopicRoomMessageCreated),
		// notification
		string(notificationDomain.TopicWelcomeEmailNotificationRequested),
		string(notificationDomain.TopicMessageNotificationRequested),
		string(notificationDomain.TopicUndeliveredMessageNotificationRequested),
	}
}

// -------------------- Event UseCases (Kafka consumer side) --------------------
func provideWebsocketUserSessionStartedUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher) usecase.UserSessionStartedUseCase {
	return usecase.NewUserSessionStartedUseCase(eventIDGen, publisher)
}
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
func provideChatPrivateMessageCreatedUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher, messageRepo *chatRepository.MessageRepository) chatUsecase.PrivateMessageCreatedUseCase {
	return chatUsecase.NewPrivateMessageCreatedUseCase(eventIDGen, messageRepo, publisher)
}
func provideChatRoomMessageCreatedUseCase(eventIDGen *uuid.EventIDGenerator, publisher *kafkautil.EventPublisher, messageRepo *chatRepository.MessageRepository) chatUsecase.RoomMessageCreatedUseCase {
	return chatUsecase.NewRoomMessageCreatedUseCase(eventIDGen, messageRepo, publisher)
}
func provideNotificationWelcomeEmailNotificationRequestedUseCase(emailNotifier *gomailUtil.EmailNotifier) notificationUsecase.WelcomeEmailNotificationRequestedUseCase {
	return notificationUsecase.NewWelcomeEmailNotificationRequestedUseCase(emailNotifier)
}
func provideNotificationMessageNotificationRequestedUseCase(messageRepo *notificationRepository.MessageRepository, messageNotifier *notificationWebsocketInfrastructure.MessageNotifier) notificationUsecase.MessageNotificationRequestedUseCase {
	return notificationUsecase.NewMessageNotificationRequestedUseCase(messageRepo, messageNotifier, messageRepo)
}
func provideNotificationUndeliveredMessageNotificationRequestedUseCase(messageRepo *notificationRepository.MessageRepository, messageNotifier *notificationWebsocketInfrastructure.MessageNotifier) notificationUsecase.UndeliveredMessageNotificationRequestedUseCase {
	return notificationUsecase.NewUndeliveredMessageNotificationRequestedUseCase(messageRepo, messageNotifier, messageRepo)
}
func provideNotificationMessageReadUseCase(messageRepo *notificationRepository.MessageRepository) notificationUsecase.MessageReadUseCase {
	return notificationUsecase.NewMessageReadUseCase(messageRepo)
}

// -------------------- http router --------------------
func provideHttpRouter(
	appConfig *config.App,
	signUp authorizationUsecase.SignUpUseCase,
	login authorizationUsecase.LoginUseCase,
	refreshAccessToken authorizationUsecase.RefreshAccessTokenUseCase,
	parseAccessToken authorizationUsecase.ParseAccessTokenUseCase,
	sendPrivateMessage chatUsecase.SendPrivateMessageUseCase,
	sendRoomMessage chatUsecase.SendRoomMessageUseCase,
	createRoom socialUseCase.CreateRoomUseCase,
	joinRoom socialUseCase.JoinRoomUseCase,
	leaveRoom socialUseCase.LeaveRoomUseCase,
	validator *ginutils.Validator,
	redisClient *redis.Client,
	websocketServer *websocket.Server,
) *gin.Engine {
	router := gin.New()

	router.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewCORSMiddleware(appConfig.CORS),
		middleware.NewRateLimitMiddleware(redisClient, appConfig.RateLimit),
	)

	router.Any("/health_check", handler.NewHealthCheckHandler())

	baseGroup := router.Group("/api/v1")

	authorizationGroup := baseGroup.Group("/authorization")

	authorizationGroup.Use()
	{
		authorizationGroup.POST("/sign_up", authorizationHttp.NewSignUpHandler(signUp, validator))
		authorizationGroup.POST("/login", authorizationHttp.NewLoginHandler(login, validator, appConfig.Cookie))
		authorizationGroup.GET("/refresh_access_token", authorizationHttp.NewRefreshAccessTokenHandler(refreshAccessToken, validator, appConfig.Cookie))
	}

	authorizationMiddleware := authorizationHttp.NewAuthorizationMiddleware(parseAccessToken)

	chatGroup := baseGroup.Group("/chat")
	chatGroup.Use(authorizationMiddleware)
	{
		chatGroup.POST("/private", chatHttp.NewSendPrivateMessageHandler(sendPrivateMessage, validator))
		chatGroup.POST("/room", chatHttp.NewSendRoomMessageHandler(sendRoomMessage, validator))
	}

	socialGroup := baseGroup.Group("/social")
	socialGroup.Use(authorizationMiddleware)
	{
		socialGroup.POST("/room", socialHttp.NewCreateRoomHandler(createRoom, validator))
		socialGroup.POST("/room/member", socialHttp.NewJoinRoomHandler(joinRoom, validator))
		socialGroup.DELETE("/room/member", socialHttp.NewLeaveRoomHandler(leaveRoom, validator))
	}

	websocketGroup := baseGroup.Group("/ws")
	websocketGroup.Use(authorizationMiddleware)
	{
		websocketGroup.GET("/", handler.NewWebsocketHandler(websocketServer))
	}

	return router
}

// -------------------- websocket router --------------------
func provideWebsocketRouter(
	notificationMessageRead notificationUsecase.MessageReadUseCase,
) *websocket.Router {
	router := websocket.NewRouter()

	router.Handle(notificationWebsocket.MessageReadRequestTopic, notificationWebsocket.NewMessageReadHandler(notificationMessageRead))

	return router
}

// -------------------- Kafka Subscriptions --------------------
func provideKafkaSubscriptions(subscriber *kafkautil.EventSubscriber,
	//websocket
	userSessionStartedUseCase usecase.UserSessionStartedUseCase,
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
	chatPrivateMessageCreated chatUsecase.PrivateMessageCreatedUseCase,
	chatRoomMessageCreated chatUsecase.RoomMessageCreatedUseCase,
	// notification
	notificationWelcomeEmailNotificationRequested notificationUsecase.WelcomeEmailNotificationRequestedUseCase,
	notificationMessageNotificationRequested notificationUsecase.MessageNotificationRequestedUseCase,
	notificationUndeliveredMessageNotificationRequested notificationUsecase.UndeliveredMessageNotificationRequestedUseCase,
) error {
	//websocket
	subscriber.Subscribe(websocket.TopicUserSessionStarted, kafka.NewUserSessionStartedEventHandler(userSessionStartedUseCase))
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
	subscriber.Subscribe(chatDomain.TopicPrivateMessageCreated, chatKafka.NewPrivateMessageCreatedEventHandler(chatPrivateMessageCreated))
	subscriber.Subscribe(chatDomain.TopicRoomMessageCreated, chatKafka.NewRoomMessageCreatedEventHandler(chatRoomMessageCreated))
	// Notification
	subscriber.Subscribe(notificationDomain.TopicWelcomeEmailNotificationRequested, notificationKafka.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested))
	subscriber.Subscribe(notificationDomain.TopicMessageNotificationRequested, notificationKafka.NewMessageNotificationRequestedEventHandler(notificationMessageNotificationRequested))
	subscriber.Subscribe(notificationDomain.TopicUndeliveredMessageNotificationRequested, notificationKafka.NewUndeliveredMessageNotificationRequestedEventHandler(notificationUndeliveredMessageNotificationRequested))
	return nil
}

func BuildDependencies(
	needMigrate bool,
	appConfig *config.App,
	ginEngine *gin.Engine,
	mysql *gorm.DB,
	topics []string,
	kafkaPublisher *kafkautil.EventPublisher,
	kafkaSubscriber *kafkautil.EventSubscriber,
	kafkaConsumerRetrier *kafkautil.ConsumerRetrier,
	kafkaProducerWithRetry *kafkautil.ProducerWithRetry,
	outboxConsumer *canal.OutboxConsumer,
	emailNotifier *gomailUtil.EmailNotifier,
	// subscriptions side-effect
	_ error, // ensure subscriptions provider executed (ignored)
) (*Dependencies, error) {
	// Migrations (side-effect). Performed here to keep initialize logic centralized.
	if needMigrate {
		if err := gormutils.AutoMigrate(
			mysql,
			&authorizationModel.User{},
			&notificationModel.Message{},
			&notificationModel.MessageState{},
			&chatModel.User{},
			&chatModel.Room{},
			&chatModel.PrivateMessage{},
			&chatModel.RoomMessage{},
			&socialModel.User{},
			&socialModel.Room{},
			&model.Event{},
			&model.DeadLetter{},
		); err != nil {
			return nil, fmt.Errorf("mysql migrate failed, err:%w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := kafkautil.EnsureTopics(ctx, appConfig.Kafka.Common, topics, 1, 1); err != nil {
		return nil, fmt.Errorf("ensure kafka topics failed, err:%w", err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	if err := kafkaSubscriber.Start(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("start kafka subscriber failed, err:%w", err)
	}
	outboxConsumer.Start()
	emailNotifier.Start()
	kafkaPublisher.Start()
	kafkaConsumerRetrier.Start()

	deps := &Dependencies{
		HttpRouter: ginEngine,
		config:     appConfig,
		closeAll: func() {
			outboxConsumer.Close()
			emailNotifier.Close()
			kafkaPublisher.Close()
			kafkaSubscriber.Close()
			kafkaConsumerRetrier.Close()
			kafkaProducerWithRetry.Close()
			cancel()
		},
	}
	return deps, nil
}
