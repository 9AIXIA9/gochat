//go:build wireinject
// +build wireinject

package main

import (
	"gochat/api"
	"gochat/internal/config"
	"gochat/internal/infra/logger"
	"gochat/internal/infra/snowflake"
	"gochat/internal/usecase"

	"github.com/google/wire"
)

// ProvideConfig 提供配置依赖
func ProvideConfig(configPath string) *config.Config {
	// 加载配置
	conf := config.MustLoad(configPath)

	// 初始化系统组件，失败直接panic
	logger.MustInit(conf.Log)
	snowflake.MustInit(conf.Snowflake)

	return conf
}

// UsecaseSet 提供所有的 Usecase
var UsecaseSet = wire.NewSet(
	usecase.NewAuth,
	usecase.NewSignup,
	usecase.NewLogin,
	usecase.NewCreateRoom,
	usecase.NewJoinRoom,
	usecase.NewExitRoom,
)

// InitializeDependencies 使用 Wire 初始化所有依赖
func InitializeDependencies(configPath string) *api.Dependencies {
	wire.Build(
		UsecaseSet,
		ProvideConfig,
		wire.Struct(new(api.Dependencies), "*"),
	)
	return &api.Dependencies{}
}
