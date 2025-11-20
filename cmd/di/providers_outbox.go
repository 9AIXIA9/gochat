package di

import (
	"gochat/config"
	canalUtil "gochat/internal/infrastructure/canal"
	kafkautil "gochat/internal/infrastructure/kafka"
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
	return canalUtil.NewCanal(appConfig.BinlogReader)
}

func provideCanalOutboxConsumer(c *canal.Canal, publisher *kafkautil.EventPublisher, eventRepo *repository.EventRepository, metrics *prometheus.Metrics) *canalUtil.OutboxConsumer {
	oc := canalUtil.NewOutboxConsumer(c, publisher, eventRepo)
	oc.SetMetrics(metrics)
	return oc
}
