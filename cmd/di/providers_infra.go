package di

import (
	"context"
	"gochat/config"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authUUID "gochat/internal/authorization/infrastructure/uuid"
	chatDomain "gochat/internal/chat/domain"
	chatUUID "gochat/internal/chat/infrastructure/uuid"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipUUID "gochat/internal/friendship/infrastructure/uuid"
	"gochat/internal/infrastructure/bcrypt"
	gormInfra "gochat/internal/infrastructure/gorm"
	kafkautil "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	redisInfra "gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/uuid"
	validatorInfra "gochat/internal/infrastructure/validator"
	"gochat/internal/infrastructure/websocket"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	gomailInfra "gochat/internal/notification/infrastructure/gomail"
	notificationWebsocket "gochat/internal/notification/infrastructure/websocket"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipSnowflake "gochat/internal/roomship/infrastructure/snowflake"
	roomshipUUID "gochat/internal/roomship/infrastructure/uuid"
	"gochat/internal/shared/event"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type emailServiceAvailable bool

var InfraSet = wire.NewSet(
	provideMysql,
	provideRedis,
	provideKafkaPublisher,
	provideValidator,
	provideMetrics,
	// Generators & managers (concrete providers)
	provideEventIDGenerator,
	provideAuthorizationUserIDGenerator,
	provideAuthorizationUserNumberGenerator,
	provideHasher,
	provideMessageIDGenerator,
	provideRoomshipRoomIDGenerator,
	provideRoomshipRoomNumberGenerator,
	provideFriendshipOperationIDGenerator,
	provideAccessTokenManager,
	provideRefreshTokenGenerator,
	provideGomailDialer,
	provideEmailAvailable,
	provideEmailNotifier,
	providePrivateMessageNotifier,
	provideRoomMessageNotifier,
	// Binds
	wire.Bind(new(event.IDGenerator), new(*uuid.EventIDGenerator)),
	wire.Bind(new(event.Publisher), new(*kafkautil.EventPublisher)),
	// Authorization binds
	wire.Bind(new(authDomain.UserIDGenerator), new(*authUUID.UserIDGenerator)),
	wire.Bind(new(authDomain.UserNumberGenerator), new(*authSnowflake.UserNumberGenerator)),
	wire.Bind(new(authDomain.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(authDomain.Comparator), new(*bcrypt.Hasher)),
	wire.Bind(new(authDomain.AccessTokenGenerator), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authApp.AccessTokenParser), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authDomain.RefreshTokenGenerator), new(*crypto.RefreshTokenGenerator)),
	// Chat binds
	wire.Bind(new(chatDomain.MessageIDGenerator), new(*chatUUID.MessageIDGenerator)),
	// Roomship generator & crypto binds
	wire.Bind(new(roomshipDomain.RoomIDGenerator), new(*roomshipUUID.RoomIDGenerator)),
	wire.Bind(new(roomshipDomain.RoomNumberGenerator), new(*roomshipSnowflake.RoomNumberGenerator)),
	wire.Bind(new(roomshipDomain.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(roomshipDomain.Comparator), new(*bcrypt.Hasher)),
	// Friendship generator
	wire.Bind(new(friendshipDomain.OperationIDGenerator), new(*friendshipUUID.OperationIDGenerator)),
	// Notification binds
	wire.Bind(new(notificationApp.WelcomeEmailNotifier), new(*gomailInfra.EmailNotifier)),
	wire.Bind(new(notificationDomain.PrivateMessageNotifier), new(*notificationWebsocket.PrivateMessageNotifier)),
	wire.Bind(new(notificationDomain.RoomMessageNotifier), new(*notificationWebsocket.RoomMessageNotifier)),
)

func provideMysql(appConfig *config.App) (*gorm.DB, error) {
	return gormInfra.ConnectToMysql(appConfig.Mysql)
}
func provideRedis(appConfig *config.App) (*redis.Client, error) {
	return redisInfra.ConnectToRedis(appConfig.Redis)
}

func provideValidator() (*validatorInfra.Validator, error) { return validatorInfra.NewValidator() }

func provideMetrics() *prometheus.Metrics { return prometheus.NewMetrics(nil) }

func provideEventIDGenerator() *uuid.EventIDGenerator { return uuid.NewEventIDGenerator() }

func provideAuthorizationUserIDGenerator() *authUUID.UserIDGenerator {
	return authUUID.NewUserIDGenerator()
}
func provideAuthorizationUserNumberGenerator(appConfig *config.App) (*authSnowflake.UserNumberGenerator, error) {
	return authSnowflake.NewUserNumberGenerator(appConfig.MachineNode)
}
func provideHasher(appConfig *config.App) *bcrypt.Hasher { return bcrypt.NewHasher(appConfig.Hasher) }
func provideMessageIDGenerator() *chatUUID.MessageIDGenerator {
	return chatUUID.NewMessageIDGenerator()
}

func provideRoomshipRoomIDGenerator() *roomshipUUID.RoomIDGenerator {
	return roomshipUUID.NewRoomIDGenerator()
}
func provideRoomshipRoomNumberGenerator(appConfig *config.App) (*roomshipSnowflake.RoomNumberGenerator, error) {
	return roomshipSnowflake.NewRoomNumberGenerator(appConfig.MachineNode)
}
func provideFriendshipOperationIDGenerator() *friendshipUUID.OperationIDGenerator {
	return friendshipUUID.NewOperationIDGenerator()
}
func provideAccessTokenManager(appConfig *config.App) *jwt.AccessTokenManager {
	return jwt.NewAccessTokenManager(appConfig.AccessToken)
}
func provideRefreshTokenGenerator(appConfig *config.App) *crypto.RefreshTokenGenerator {
	return crypto.NewRefreshTokenGenerator(appConfig.RefreshToken)
}
func provideGomailDialer(appConfig *config.App) (*gomail.Dialer, error) {
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
func providePrivateMessageNotifier(manager *websocket.Manager) *notificationWebsocket.PrivateMessageNotifier {
	return notificationWebsocket.NewPrivateMessageNotifier(manager)
}
func provideRoomMessageNotifier(manager *websocket.Manager) *notificationWebsocket.RoomMessageNotifier {
	return notificationWebsocket.NewRoomMessageNotifier(manager)
}
func provideKafkaPublisher(appConfig *config.App, eventRepo *repository.EventRepository, metrics *prometheus.Metrics) (*kafkautil.EventPublisher, error) {
	return kafkautil.NewEventPublisher(appConfig.Kafka, metrics, func(id event.ID) error {
		return eventRepo.MarkAsPublished(context.Background(), id)
	})
}
