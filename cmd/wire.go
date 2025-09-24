//go:build wireinject
// +build wireinject

package main

import (
	"github.com/redis/go-redis/v9"
	"gochat/api"
	"gochat/internal/config"
	"gochat/internal/handler"
	"gochat/internal/infra/logger"
	"gochat/internal/infra/repository"
	"gochat/internal/infra/snowflake"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/usecase"
	"gorm.io/gorm"

	"github.com/google/wire"
)

// ProvideConfig 提供配置依赖
func ProvideConfig(configPath string) *config.Config {
	// 加载配置
	conf := config.MustLoad(configPath)

	conf.Validate()

	// 初始化系统组件，失败直接panic
	logger.MustInit(conf.Log)
	snowflake.MustInit(conf.Snowflake)
	handler.MustInitTrans(conf.Language)

	return conf
}

// ProvideMysqlConnection 提供 Mysql数据库连接
func ProvideMysqlConnection(conf *config.Config) *gorm.DB {
	return repository.MustConnectToMysql(conf.Database)
}

// ProvideRedisConnection 提供 Redis数据库连接
func ProvideRedisConnection(conf *config.Config) *redis.Client {
	return repository.MustConnectToRedis(conf.Redis)
}

func ProvideTokenConf(conf *config.Config) *config.Token {
	return conf.Token
}

func ProvideRefreshTokenConf(conf *config.Config) *config.RefreshToken {
	return conf.Token.Refresh
}

// RepositorySet 提供所有的Repository
var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewRoomRepository,
	repository.NewMessageRepository,
	repository.NewUserRoomRepository,
	repository.NewRefreshTokenRepository,
)

// UsecaseSet 提供所有的Usecase
var UsecaseSet = wire.NewSet(
	// Auth
	usecase.NewAuth,

	// Login
	usecase.NewLogin,

	// Signup
	usecase.NewSignup,

	// RefreshToken
	usecase.NewRefreshToken,

	// Room相关
	usecase.NewCreateRoom,
	usecase.NewJoinRoom,
	usecase.NewLeaveRoom,

	// message相关
	usecase.NewSendMessage,

	// websocket相关
	usecase.NewUserConnected,
)

// InitializeDependencies 使用Wire初始化所有依赖
func InitializeDependencies(configPath string) *api.Dependencies {
	wire.Build(
		ProvideConfig,
		ProvideTokenConf,
		ProvideRefreshTokenConf,
		ProvideMysqlConnection,
		ProvideRedisConnection,
		RepositorySet,
		UsecaseSet,
		manager.NewManager,
		wire.Struct(new(api.Dependencies), "*"),
	)
	return &api.Dependencies{}
}
