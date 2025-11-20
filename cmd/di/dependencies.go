package di

import (
	canalUtil "gochat/internal/infrastructure/canal"
	ginutils "gochat/internal/infrastructure/gin"
	kafkautil "gochat/internal/infrastructure/kafka"
	gomailUtil "gochat/internal/notification/infrastructure/gomail"
)

type Dependencies struct {
	HttpServer           *ginutils.Server
	KafkaEventPublisher  *kafkautil.EventPublisher
	KafkaEventSubscriber *kafkautil.EventSubscriber
	OutboxConsumer       *canalUtil.OutboxConsumer
	EmailNotifier        *gomailUtil.EmailNotifier
}

func BuildDependencies(
	httpServer *ginutils.Server,
	kafkaPublisher *kafkautil.EventPublisher,
	kafkaSubscriber *kafkautil.EventSubscriber,
	outboxConsumer *canalUtil.OutboxConsumer,
	emailNotifier *gomailUtil.EmailNotifier,
	_ emailServiceAvailable,
	_ kafkaTopicEnsured,
	_ kafkaTopicSubscribed,
	_ databaseMigrated,
) (*Dependencies, error) {
	deps := &Dependencies{
		HttpServer:           httpServer,
		KafkaEventPublisher:  kafkaPublisher,
		KafkaEventSubscriber: kafkaSubscriber,
		OutboxConsumer:       outboxConsumer,
		EmailNotifier:        emailNotifier,
	}
	return deps, nil
}
