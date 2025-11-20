package di

import (
	"gochat/config"
	authorizationApp "gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authorizationSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authorizationUuid "gochat/internal/authorization/infrastructure/uuid"
	chatApp "gochat/internal/chat/application"
	chatUUID "gochat/internal/chat/infrastructure/uuid"
	chatUuid "gochat/internal/chat/infrastructure/uuid"
	"gochat/internal/infrastructure/bcrypt"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/prometheus"
	redisutils "gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/uuid"
	ginutils "gochat/internal/infrastructure/validator"
	"gochat/internal/infrastructure/websocket"
	gomailUtil "gochat/internal/notification/infrastructure/gomail"
	notificationWebsocketInfrastructure "gochat/internal/notification/infrastructure/websocket"
	"gochat/internal/shared/event"
	socialApp "gochat/internal/social/application"
	socialSnowflake "gochat/internal/social/infrastructure/snowflake"
	socialUUID "gochat/internal/social/infrastructure/uuid"
	socialUuid "gochat/internal/social/infrastructure/uuid"

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
	provideValidator,
	provideMetrics,
	// Generators & managers (concrete providers)
	provideEventIDGenerator,
	provideAuthorizationUserIDGenerator,
	provideAuthorizationUserNumberGenerator,
	provideHasher,
	provideMessageIDGenerator,
	provideSocialRoomIDGenerator,
	provideSocialRoomNumberGenerator,
	provideAccessTokenManager,
	provideRefreshTokenGenerator,
	provideGomailDialer,
	provideEmailAvailable,
	provideEmailNotifier,
	provideMessageNotifier,
	// Binds
	wire.Bind(new(event.IDGenerator), new(*uuid.EventIDGenerator)),
	// Authorization binds
	wire.Bind(new(authorizationApp.UserIDGenerator), new(*authorizationUuid.UserIDGenerator)),
	wire.Bind(new(authorizationApp.UserNumberGenerator), new(*authorizationSnowflake.UserNumberGenerator)),
	wire.Bind(new(authorizationApp.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(authorizationApp.Comparator), new(*bcrypt.Hasher)),
	wire.Bind(new(authorizationApp.AccessTokenGenerator), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authorizationApp.AccessTokenParser), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authorizationApp.RefreshTokenGenerator), new(*crypto.RefreshTokenGenerator)),
	// Chat binds
	wire.Bind(new(chatApp.MessageIDGenerator), new(*chatUUID.MessageIDGenerator)),
	// Social generator & crypto binds
	wire.Bind(new(socialApp.RoomIDGenerator), new(*socialUUID.RoomIDGenerator)),
	wire.Bind(new(socialApp.RoomNumberGenerator), new(*socialSnowflake.RoomNumberGenerator)),
	wire.Bind(new(socialApp.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(socialApp.Comparator), new(*bcrypt.Hasher)),
)

func provideMysql(appConfig *config.App) (*gorm.DB, error) {
	return gormutils.ConnectToMysql(appConfig.Mysql)
}
func provideRedis(appConfig *config.App) (*redis.Client, error) {
	return redisutils.ConnectToRedis(appConfig.Redis)
}

func provideValidator() (*ginutils.Validator, error) { return ginutils.NewValidator() }

func provideMetrics() *prometheus.Metrics { return prometheus.NewMetrics(nil) }

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
func provideEmailAvailable(d *gomail.Dialer) emailServiceAvailable {
	if err := gomailUtil.TestConnection(d); err != nil {
		zap.L().Info("Email dialer connection test failed, email notifier will be disabled", zap.Error(err))
		return false
	}
	return true
}
func provideEmailNotifier(appConfig *config.App, dialer *gomail.Dialer) *gomailUtil.EmailNotifier {
	return gomailUtil.NewEmailNotifier(appConfig.Name, dialer)
}
func provideMessageNotifier(manager *websocket.Manager) *notificationWebsocketInfrastructure.MessageNotifier {
	return notificationWebsocketInfrastructure.NewMessageNotifier(manager)
}
