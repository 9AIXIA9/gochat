//go:build wireinject
// +build wireinject

package di

import (
	"gochat/config"

	"github.com/google/wire"
)

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

// Initialize wires up the dependencies for main
func Initialize(appConfig *config.App) (*Dependencies, error) {
	wire.Build(
		InfraSet,
		RepoSet,
		KafkaSet,
		BinlogSet,
		UseCaseHTTPSet,
		UseCaseWebsocketSet,
		UseCaseKafkaSet,
		WebsocketSet,
		HTTPSet,
		BuildDependencies,
	)
	return &Dependencies{}, nil
}
