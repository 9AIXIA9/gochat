package di

import (
	"context"
	"gochat/config"
	authDomain "gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authUUID "gochat/internal/authorization/infrastructure/uuid"
	chatDomain "gochat/internal/chat/domain"
	chatWebsocket "gochat/internal/chat/infrastructure/websocket"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipUUID "gochat/internal/friendship/infrastructure/uuid"
	"gochat/internal/infrastructure/bcrypt"
	ginutils "gochat/internal/infrastructure/gin"
	gormInfra "gochat/internal/infrastructure/gorm"
	kafkautil "gochat/internal/infrastructure/kafka"
	infraotel "gochat/internal/infrastructure/otel"
	redisInfra "gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/uuid"
	validatorInfra "gochat/internal/infrastructure/validator"
	"gochat/internal/infrastructure/websocket"
	notificationDomain "gochat/internal/notification/domain"
	gomailInfra "gochat/internal/notification/infrastructure/gomail"
	notificationWebsocket "gochat/internal/notification/infrastructure/websocket"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipSnowflake "gochat/internal/roomship/infrastructure/snowflake"
	roomshipUUID "gochat/internal/roomship/infrastructure/uuid"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type (
	emailServiceAvailable bool
	OTELShutdown          func(context.Context) error
)

var InfraSet = wire.NewSet(
	provideObservability,
	provideMysqlConnection,
	provideRedisConnection,
	provideKafkaPublisher,
	provideValidator,
	// Generators & managers (concrete providers)
	provideEventIDGenerator,
	provideAuthorizationUserIDGenerator,
	provideAuthorizationUserNumberGenerator,
	provideHasher,
	provideMessageIDGenerator,
	provideRoomshipRoomIDGenerator,
	provideRoomshipRoomshipIDGenerator,
	provideRoomshipRoomNumberGenerator,
	provideOperationIDGenerator,
	provideFriendshipFriendshipIDGenerator,
	provideAccessTokenManager,
	provideRefreshTokenGenerator,
	provideGomailDialer,
	provideEmailAvailable,
	provideEmailNotifier,
	provideSystemMessageNotifier,
	providePrivateMessageNotifier,
	provideRoomMessageNotifier,
	// Binds
	wire.Bind(new(event.IDGenerator), new(*uuid.EventIDGenerator)),
	wire.Bind(new(kernel.MessageIDGenerator), new(*uuid.MessageIDGenerator)),
	wire.Bind(new(kernel.OperationIDGenerator), new(*uuid.OperationIDGenerator)),
	wire.Bind(new(event.Publisher), new(*kafkautil.EventPublisher)),
	wire.Bind(new(ginutils.Validator), new(*validatorInfra.Validator)),
	wire.Bind(new(websocket.Validator), new(*validatorInfra.Validator)),
	// Authorization binds
	wire.Bind(new(authDomain.UserIDGenerator), new(*authUUID.UserIDGenerator)),
	wire.Bind(new(authDomain.UserNumberGenerator), new(*authSnowflake.UserNumberGenerator)),
	wire.Bind(new(authDomain.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(authDomain.Comparator), new(*bcrypt.Hasher)),
	wire.Bind(new(authDomain.AccessTokenGenerator), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authDomain.AccessTokenParser), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authDomain.RefreshTokenGenerator), new(*crypto.RefreshTokenGenerator)),
	// Chat Notifier
	wire.Bind(new(chatDomain.PrivateMessageNotifier), new(*chatWebsocket.PrivateMessageNotifier)),
	wire.Bind(new(chatDomain.RoomMessageNotifier), new(*chatWebsocket.RoomMessageNotifier)),
	// Roomship generator & crypto binds
	wire.Bind(new(roomshipDomain.RoomIDGenerator), new(*roomshipUUID.RoomIDGenerator)),
	wire.Bind(new(roomshipDomain.RoomshipIDGenerator), new(*roomshipUUID.RoomshipIDGenerator)),
	wire.Bind(new(roomshipDomain.RoomNumberGenerator), new(*roomshipSnowflake.RoomNumberGenerator)),
	wire.Bind(new(roomshipDomain.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(roomshipDomain.Comparator), new(*bcrypt.Hasher)),
	// Friendship generator
	wire.Bind(new(friendshipDomain.FriendshipIDGenerator), new(*friendshipUUID.FriendshipIDGenerator)),
	// Notification binds
	wire.Bind(new(notificationDomain.WelcomeEmailNotifier), new(*gomailInfra.EmailNotifier)),
	wire.Bind(new(notificationDomain.SystemMessageNotifier), new(*notificationWebsocket.SystemMessageNotifier)),
)

func provideObservability(cfg *config.App) (OTELShutdown, error) {
	if cfg.OTEL == nil {
		return func(context.Context) error { return nil }, nil
	}
	return infraotel.Init(cfg.OTEL)
}

func provideMysqlConnection(appConfig *config.App) (*gorm.DB, error) {
	return gormInfra.ConnectToMysql(appConfig.Mysql)
}
func provideRedisConnection(appConfig *config.App) (*redis.Client, error) {
	return redisInfra.ConnectToRedis(appConfig.Redis)
}

func provideValidator() (*validatorInfra.Validator, error) { return validatorInfra.NewValidator() }

func provideEventIDGenerator() *uuid.EventIDGenerator {
	return uuid.NewEventIDGenerator()
}
func provideOperationIDGenerator() *uuid.OperationIDGenerator {
	return uuid.NewOperationIDGenerator()
}
func provideMessageIDGenerator() *uuid.MessageIDGenerator {
	return uuid.NewMessageIDGenerator()
}

func provideAuthorizationUserIDGenerator() *authUUID.UserIDGenerator {
	return authUUID.NewUserIDGenerator()
}
func provideAuthorizationUserNumberGenerator(appConfig *config.App) (*authSnowflake.UserNumberGenerator, error) {
	return authSnowflake.NewUserNumberGenerator(appConfig.MachineNode)
}
func provideHasher(appConfig *config.App) *bcrypt.Hasher {
	return bcrypt.NewHasher(appConfig.Hasher)
}
func provideRoomshipRoomIDGenerator() *roomshipUUID.RoomIDGenerator {
	return roomshipUUID.NewRoomIDGenerator()
}
func provideRoomshipRoomshipIDGenerator() *roomshipUUID.RoomshipIDGenerator {
	return roomshipUUID.NewRoomshipIDGenerator()
}
func provideRoomshipRoomNumberGenerator(appConfig *config.App) (*roomshipSnowflake.RoomNumberGenerator, error) {
	return roomshipSnowflake.NewRoomNumberGenerator(appConfig.MachineNode)
}
func provideFriendshipFriendshipIDGenerator() *friendshipUUID.FriendshipIDGenerator {
	return friendshipUUID.NewFriendshipIDGenerator()
}
func provideAccessTokenManager(appConfig *config.App) *jwt.AccessTokenManager {
	return jwt.NewAccessTokenManager(appConfig.AccessToken)
}
func provideRefreshTokenGenerator(appConfig *config.App) *crypto.RefreshTokenGenerator {
	return crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)
}
func provideGomailDialer(appConfig *config.App) *gomail.Dialer {
	return gomailInfra.NewDialer(appConfig.Email)
}
func provideEmailAvailable(d *gomail.Dialer) emailServiceAvailable {
	if err := gomailInfra.TestConnection(d); err != nil {
		zap.L().Info("Email dialer connection test failed, email notifier will be disabled", zap.Error(err))
		return false
	}
	return true
}
func provideEmailNotifier(appConfig *config.App, dialer *gomail.Dialer) *gomailInfra.EmailNotifier {
	return gomailInfra.NewEmailNotifier(appConfig.Name, dialer)
}
func provideSystemMessageNotifier(manager *websocket.Manager) *notificationWebsocket.SystemMessageNotifier {
	return notificationWebsocket.NewSystemMessageNotifier(manager)
}
func providePrivateMessageNotifier(manager *websocket.Manager) *chatWebsocket.PrivateMessageNotifier {
	return chatWebsocket.NewPrivateMessageNotifier(manager)
}
func provideRoomMessageNotifier(manager *websocket.Manager) *chatWebsocket.RoomMessageNotifier {
	return chatWebsocket.NewRoomMessageNotifier(manager)
}
func provideKafkaPublisher(appConfig *config.App, eventRepo event.Repository) (*kafkautil.EventPublisher, error) {
	return kafkautil.NewEventPublisher(
		appConfig.Kafka,
		func(id event.ID) error {
			return eventRepo.MarkAsPublished(context.Background(), id)
		},
	)
}
