package di

import (
	"gochat/config"
	canalInfra "gochat/internal/infrastructure/canal"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/google/wire"
)

var CanalSet = wire.NewSet(
	provideCanal,
	provideCanalOutboxConsumer,
)

func provideCanal(appConfig *config.App) (*canal.Canal, error) {
	return canalInfra.NewCanal(appConfig.BinlogReader)
}

func provideCanalOutboxConsumer(c *canal.Canal, publisher *kafkaInfra.EventPublisher, eventRepo *repository.EventRepository, metrics *prometheus.Metrics) *canalInfra.OutboxConsumer {
	oc := canalInfra.NewOutboxConsumer(c, publisher, eventRepo)
	oc.SetMetrics(metrics)
	return oc
}
