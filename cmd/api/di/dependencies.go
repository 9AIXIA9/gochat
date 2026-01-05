package di

import (
	canalUtil "gochat/internal/infrastructure/canal"
	ginutils "gochat/internal/infrastructure/gin"
	kafkautil "gochat/internal/infrastructure/kafka"
	gomailUtil "gochat/internal/notification/infrastructure/gomail"
)

type Dependencies struct {
	HttpServer          *ginutils.Server
	KafkaEventPublisher *kafkautil.EventPublisher
	KafkaConsumers      []*kafkautil.Consumer
	BinlogReader        *canalUtil.BinlogReader
	EmailNotifier       *gomailUtil.EmailNotifier
	OTELShutdown        OTELShutdown
}

func BuildDependencies(
	httpServer *ginutils.Server,
	kafkaPublisher *kafkautil.EventPublisher,
	consumers []*kafkautil.Consumer,
	binlogReader *canalUtil.BinlogReader,
	emailNotifier *gomailUtil.EmailNotifier,
	OTELShutdown OTELShutdown,
	_ emailServiceAvailable,
) (*Dependencies, error) {
	deps := &Dependencies{
		HttpServer:          httpServer,
		KafkaEventPublisher: kafkaPublisher,
		KafkaConsumers:      consumers,
		BinlogReader:        binlogReader,
		EmailNotifier:       emailNotifier,
		OTELShutdown:        OTELShutdown,
	}
	return deps, nil
}
