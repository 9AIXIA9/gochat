package di

import (
	"context"
	"gochat/config"
	rootapp "gochat/internal/application"
	authDomain "gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	authSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authUUID "gochat/internal/authorization/infrastructure/uuid"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipUUID "gochat/internal/friendship/infrastructure/uuid"
	gatewayUUID "gochat/internal/gateway/infrastructure/uuid"
	"gochat/internal/infrastructure/bcrypt"
	ginutils "gochat/internal/infrastructure/gin"
	gormInfra "gochat/internal/infrastructure/gorm"
	kafkautil "gochat/internal/infrastructure/kafka"
	infraotel "gochat/internal/infrastructure/otel"
	outboxUtil "gochat/internal/infrastructure/outbox"
	redisInfra "gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/ulule"
	"gochat/internal/infrastructure/uuid"
	validatorInfra "gochat/internal/infrastructure/validator"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipSnowflake "gochat/internal/roomship/infrastructure/snowflake"
	roomshipUUID "gochat/internal/roomship/infrastructure/uuid"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/iputil"
	"net/http"
	"strings"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	OTELShutdown func(context.Context) error
	HTTPLimiter  limiter.Limiter
	KafkaLimiter limiter.Limiter
	checkOrigin  func(r *http.Request) bool
)

var InfraSet = wire.NewSet(
	provideObservability,
	provideMysqlConnection,
	provideRedisConnection,
	provideKafkaAsyncPublisher,
	provideKafkaSyncPublisher,
	provideOutboxDispatcher,
	provideCheckOrigin,
	provideValidator,
	provideHTTPLimiter,
	provideKafkaLimiter,
	// Generators & managers (concrete providers)
	provideEventIDGenerator,
	provideAuthorizationUserIDGenerator,
	provideAuthorizationUserNumberGenerator,
	provideHasher,
	provideMessageIDGenerator,
	provideSessionIDGenerator,
	provideRoomshipRoomIDGenerator,
	provideRoomshipRoomshipIDGenerator,
	provideRoomshipRoomNumberGenerator,
	provideOperationIDGenerator,
	provideFriendshipFriendshipIDGenerator,
	provideAccessTokenManager,
	provideRefreshTokenGenerator,
	// Binds
	wire.Bind(new(event.IDGenerator), new(*uuid.EventIDGenerator)),
	wire.Bind(new(kernel.MessageIDGenerator), new(*uuid.MessageIDGenerator)),
	wire.Bind(new(kernel.OperationIDGenerator), new(*uuid.OperationIDGenerator)),
	wire.Bind(new(event.AsyncPublisher), new(*kafkautil.EventAsyncPublisher)),
	wire.Bind(new(event.SyncPublisher), new(*kafkautil.EventSyncPublisher)),
	wire.Bind(new(ginutils.Validator), new(*validatorInfra.Validator)),
	// Authorization binds
	wire.Bind(new(authDomain.UserIDGenerator), new(*authUUID.UserIDGenerator)),
	wire.Bind(new(authDomain.UserNumberGenerator), new(*authSnowflake.UserNumberGenerator)),
	wire.Bind(new(authDomain.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(authDomain.Comparator), new(*bcrypt.Hasher)),
	wire.Bind(new(authDomain.AccessTokenGenerator), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authDomain.AccessTokenParser), new(*jwt.AccessTokenManager)),
	wire.Bind(new(authDomain.RefreshTokenGenerator), new(*crypto.RefreshTokenGenerator)),
	// Roomship generator & crypto binds
	wire.Bind(new(roomshipDomain.RoomIDGenerator), new(*roomshipUUID.RoomIDGenerator)),
	wire.Bind(new(roomshipDomain.RoomshipIDGenerator), new(*roomshipUUID.RoomshipIDGenerator)),
	wire.Bind(new(roomshipDomain.RoomNumberGenerator), new(*roomshipSnowflake.RoomNumberGenerator)),
	wire.Bind(new(roomshipDomain.Encryptor), new(*bcrypt.Hasher)),
	wire.Bind(new(roomshipDomain.Comparator), new(*bcrypt.Hasher)),
	// Friendship generator
	wire.Bind(new(friendshipDomain.FriendshipIDGenerator), new(*friendshipUUID.FriendshipIDGenerator)),
)

func provideObservability(cfg *config.App) (OTELShutdown, error) {
	if cfg.OTEL == nil || !cfg.OTEL.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	shutdown, err := infraotel.Init(cfg.OTEL)
	if err != nil {
		zap.L().Warn("OTEL disabled because initialization failed; check OTEL endpoint configuration and ensure the collector is running", zap.Error(err))
		return func(context.Context) error { return nil }, nil
	}
	return shutdown, nil
}

func provideMysqlConnection(appConfig *config.App) (*gorm.DB, error) {
	return gormInfra.ConnectToMysql(appConfig.Mysql, appConfig.OTEL != nil && appConfig.OTEL.Enabled)
}
func provideRedisConnection(appConfig *config.App) (*redis.Client, error) {
	return redisInfra.ConnectToRedis(appConfig.Redis, appConfig.OTEL != nil && appConfig.OTEL.Enabled)
}

func provideValidator() (*validatorInfra.Validator, error) { return validatorInfra.NewValidator() }

func provideHTTPLimiter(client *redis.Client, conf *config.App) *HTTPLimiter {
	return (*HTTPLimiter)(ulule.NewLimiter(client, conf.HTTPRateLimit))
}
func provideKafkaLimiter(client *redis.Client, conf *config.App) *KafkaLimiter {
	return (*KafkaLimiter)(ulule.NewLimiter(client, conf.KafkaRateLimit))
}

func provideEventIDGenerator() *uuid.EventIDGenerator {
	return uuid.NewEventIDGenerator()
}
func provideOperationIDGenerator() *uuid.OperationIDGenerator {
	return uuid.NewOperationIDGenerator()
}
func provideMessageIDGenerator() *uuid.MessageIDGenerator {
	return uuid.NewMessageIDGenerator()
}

func provideSessionIDGenerator() *gatewayUUID.SessionIDGenerator {
	return gatewayUUID.NewSessionIDGenerator()
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
func provideKafkaAsyncPublisher(appConfig *config.App, eventRepo event.Repository) (*kafkautil.EventAsyncPublisher, error) {
	return kafkautil.NewEventAsyncPublisher(
		appConfig.Kafka,
		func(id event.ID) error {
			return eventRepo.MarkAsPublished(context.Background(), id)
		},
	)
}

func provideKafkaSyncPublisher(appConfig *config.App) (*kafkautil.EventSyncPublisher, error) {
	return kafkautil.NewEventSyncPublisher(appConfig.Kafka)
}

func provideOutboxDispatcher(appConfig *config.App, uc rootapp.UnpublishedEventsCreatedUseCase) (*outboxUtil.Dispatcher, error) {
	if appConfig == nil || appConfig.Outbox == nil {
		return outboxUtil.NewDispatcher(uc, time.Second, 1)
	}
	return outboxUtil.NewDispatcher(uc, appConfig.Outbox.SweepInterval, appConfig.Outbox.TriggerBuffer)
}

func provideCheckOrigin(appConfig *config.App) checkOrigin {
	allowed := make(map[string]struct{}, len(appConfig.CORS.AllowOrigins))
	allowAll := false
	for _, o := range appConfig.CORS.AllowOrigins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			allowAll = true
			continue
		}
		allowed[o] = struct{}{}
	}

	return func(r *http.Request) bool {
		// 若未配置 origins 或包含通配符则允许所有来源
		if allowAll || len(allowed) == 0 {
			return true
		}

		if iputil.IsLoopbackRequest(r.RemoteAddr) {
			return true
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}

		_, ok := allowed[origin]
		return ok
	}
}
