//go:build wireinject
// +build wireinject

package main

import (
	"gochat/api"
	"gochat/internal/config"
	"gochat/internal/domain"
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

// ProvideDatabase 提供数据库连接
func ProvideDatabase(conf *config.Config) *gorm.DB {
	return repository.MustConnectToMysql(*conf.Database)
}

// ProvideAuth 提供Auth服务并注入JWT配置
func ProvideAuth(conf *config.Config) domain.AuthUsecase {
	return usecase.NewAuth(conf.JWT)
}

// ProvideLogin 提供Login服务并注入JWT配置
func ProvideLogin(repo domain.UserRepository, conf *config.Config) domain.LoginUsecase {
	return usecase.NewLogin(conf.JWT, repo)
}

// RepositorySet 提供所有的Repository
var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewRoomRepository,
	repository.NewMessageRepository,
	repository.NewUserRoomRepository,
)

// UsecaseSet 提供所有的Usecase
var UsecaseSet = wire.NewSet(
	// Auth
	ProvideAuth,

	// Execute
	ProvideLogin,

	// Signup
	usecase.NewSignup,

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
		ProvideDatabase,
		RepositorySet,
		UsecaseSet,
		manager.NewManager,
		wire.Struct(new(api.Dependencies), "*"),
	)
	return &api.Dependencies{}
}
