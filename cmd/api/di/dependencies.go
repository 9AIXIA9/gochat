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
	// per-context consumers (distinct types for Wire)
	authConsumer AuthKafkaConsumer,
	profileConsumer ProfileKafkaConsumer,
	chatConsumer ChatKafkaConsumer,
	notificationConsumer NotificationKafkaConsumer,
	roomshipConsumer RoomshipKafkaConsumer,
	friendshipConsumer FriendshipKafkaConsumer,
	binlogReader *canalUtil.BinlogReader,
	emailNotifier *gomailUtil.EmailNotifier,
	OTELShutdown OTELShutdown,
	_ emailServiceAvailable,
	_ databaseMigrated,
) (*Dependencies, error) {
	deps := &Dependencies{
		HttpServer:          httpServer,
		KafkaEventPublisher: kafkaPublisher,
		KafkaConsumers: []*kafkautil.Consumer{
			authConsumer,
			profileConsumer,
			chatConsumer,
			notificationConsumer,
			roomshipConsumer,
			friendshipConsumer,
		},
		BinlogReader:  binlogReader,
		EmailNotifier: emailNotifier,
		OTELShutdown:  OTELShutdown,
	}
	return deps, nil
}
