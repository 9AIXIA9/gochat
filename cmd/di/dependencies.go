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
	KafkaConsumer       *kafkautil.Consumer
	BinlogReader        *canalUtil.BinlogReader
	EmailNotifier       *gomailUtil.EmailNotifier
}

func BuildDependencies(
	httpServer *ginutils.Server,
	kafkaPublisher *kafkautil.EventPublisher,
	kafkaConsumer *kafkautil.Consumer,
	binlogReader *canalUtil.BinlogReader,
	emailNotifier *gomailUtil.EmailNotifier,
	_ emailServiceAvailable,
	_ kafkaTopicEnsured,
	_ databaseMigrated,
) (*Dependencies, error) {
	deps := &Dependencies{
		HttpServer:          httpServer,
		KafkaEventPublisher: kafkaPublisher,
		KafkaConsumer:       kafkaConsumer,
		BinlogReader:        binlogReader,
		EmailNotifier:       emailNotifier,
	}
	return deps, nil
}
