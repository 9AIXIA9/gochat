//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gochat/api"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/handler"
	"gochat/internal/infra/crypto"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/jwt"
	"gochat/internal/infra/logger"
	"gochat/internal/infra/repository"
	"gochat/internal/infra/snowflake"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/infra/websocket/upgrader"
	"gochat/internal/usecase"
	"gorm.io/gorm"
	"time"
)

// loadConfig loads application configuration from the given file path.
func loadConfig(path string) (*config.Config, error) { return config.Load(path) }

// initialize runs logger and translator initialization as a provider for Wire.
// It returns a dummy token type to participate in the graph.
type appInitialized struct{}

func initialize(cfg *config.Config) (*appInitialized, error) {
	if err := logger.Init(cfg.Log); err != nil {
		return nil, err
	}
	if err := handler.InitTrans(cfg.Language); err != nil {
		return nil, err
	}
	return &appInitialized{}, nil
}

// connectMySQL connects to MySQL using the configuration.
func connectMySQL(cfg *config.Config) (*gorm.DB, error) {
	return repository.ConnectToMysql(cfg.Database)
}

// connectRedis connects to Redis using the configuration.
func connectRedis(cfg *config.Config) (*redis.Client, error) {
	return repository.ConnectToRedis(cfg.Redis)
}

// newRepo constructs the aggregate repository from DB and Redis.
func newRepo(db *gorm.DB, rdb *redis.Client) *repository.Repository {
	return repository.NewRepository(db, rdb)
}

// Narrow repository interface providers to help Wire resolve dependencies.

// Standalone interface
func userSaver(r *repository.Repository) domain.UserSaver   { return r.Users() }
func userFinder(r *repository.Repository) domain.UserFinder { return r.Users() }

func roomSaver(r *repository.Repository) domain.RoomSaver               { return r.Rooms() }
func roomFinder(r *repository.Repository) domain.RoomFinder             { return r.Rooms() }
func roomJoiner(r *repository.Repository) domain.RoomJoiner             { return r.Rooms() }
func roomLeaver(r *repository.Repository) domain.RoomLeaver             { return r.Rooms() }
func roomMemberFinder(r *repository.Repository) domain.RoomMemberFinder { return r.Rooms() }

func messageSaver(r *repository.Repository) domain.MessageSaver               { return r.Messages() }
func unsentMessageFinder(r *repository.Repository) domain.UnsentMessageFinder { return r.Messages() }
func sentMessageUpdater(r *repository.Repository) domain.SentMessageUpdater   { return r.Messages() }

func refreshTokenSaver(r *repository.Repository) domain.RefreshTokenSaver   { return r.RefreshTokens() }
func refreshTokenFinder(r *repository.Repository) domain.RefreshTokenFinder { return r.RefreshTokens() }

// Aggregation interface

func joinRoomAggregateUOW(r *repository.Repository) domain.JoinRoomAggregateUOW {
	return r.Rooms()
}

func sendMessageAggregate(r *repository.Repository) domain.SendMessageAggregate { return r.Messages() }

func refreshTokenAggregate(r *repository.Repository) domain.RefreshTokenAggregate {
	return r.RefreshTokens()
}

func sendUnsentMessageAggregate(r *repository.Repository) domain.SendUnsentMessageAggregate {
	return r.Messages()
}

func newWebsocketManager() *manager.Manager { return manager.NewWebsocketManager() }
func newUpgrader(cfg *config.Config) *upgrader.Upgrader {
	return upgrader.NewUpgrader(cfg.CORS.Origins)
}

func newAuthTokenManager(cfg *config.Config) *jwt.AuthTokenManager {
	return jwt.NewAuthTokenManager(cfg.Token.Auth)
}
func newRefreshTokenGenerator(cfg *config.Config) (*crypto.RefreshTokenGenerator, error) {
	return crypto.NewRefreshTokenGenerator(cfg.Token.Refresh.Length)
}
func newNumberGenerator(cfg *config.Config) (*snowflake.NumberGenerator, error) {
	return snowflake.NewNumberGenerator(cfg.Snowflake)
}
func newEncryptor(cfg *config.Config) *encrypt.BcryptEncryptor {
	return encrypt.NewBcryptEncryptor(cfg.Encryptor)
}
func refreshExpire(cfg *config.Config) time.Duration { return cfg.Token.Refresh.ExpireDuration }

// provider set for configuration and side-effects
var configSet = wire.NewSet(
	loadConfig,
	initialize,
)

// provider set for infrastructure services
var infraSet = wire.NewSet(
	connectMySQL,
	connectRedis,
	newRepo,
	newWebsocketManager,
	newUpgrader,
	newAuthTokenManager,
	newRefreshTokenGenerator,
	newNumberGenerator,
	newEncryptor,
	// repo interfaces
	userSaver,
	userFinder,
	roomSaver,
	roomFinder,
	roomJoiner,
	roomLeaver,
	roomMemberFinder,
	joinRoomAggregateUOW,
	messageSaver,
	unsentMessageFinder,
	sentMessageUpdater,
	sendMessageAggregate,
	sendUnsentMessageAggregate,
	refreshTokenSaver,
	refreshTokenFinder,
	refreshTokenAggregate,
	refreshExpire,
	// Bind concrete implementations to domain interfaces
	wire.Bind(new(domain.AuthTokenParser), new(*jwt.AuthTokenManager)),
	wire.Bind(new(domain.AuthTokenGenerator), new(*jwt.AuthTokenManager)),
	wire.Bind(new(domain.Encryptor), new(*encrypt.BcryptEncryptor)),
	wire.Bind(new(domain.Comparator), new(*encrypt.BcryptEncryptor)),
	wire.Bind(new(domain.NumberGenerator), new(*snowflake.NumberGenerator)),
	wire.Bind(new(domain.MessageSender), new(*manager.Manager)),
	wire.Bind(new(domain.RefreshTokenGenerator), new(*crypto.RefreshTokenGenerator)),
)

// provider set for usecases
var usecaseSet = wire.NewSet(
	usecase.NewAuth,
	usecase.NewSignup,
	usecase.NewLogin,
	usecase.NewRefreshToken,
	usecase.NewCreateRoom,
	usecase.NewJoinRoom,
	usecase.NewLeaveRoom,
	usecase.NewSendPrivateMessage,
	usecase.NewSendRoomMessage,
	usecase.NewUserConnected,
)

// NewDependencies is the Wire injector that builds all dependencies for the API server.
func NewDependencies(configPath string) (*api.Dependencies, error) {
	wire.Build(
		configSet,
		infraSet,
		usecaseSet,
		// Fill the Dependencies struct fields from the provided types
		wire.Struct(new(api.Dependencies), "*"),
	)
	return nil, nil
}
