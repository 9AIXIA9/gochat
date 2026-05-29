package di

import (
	canalUtil "gochat/internal/infrastructure/canal"
	ginutils "gochat/internal/infrastructure/gin"
	kafkautil "gochat/internal/infrastructure/kafka"
	outboxUtil "gochat/internal/infrastructure/outbox"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	HttpServer               *ginutils.Server
	MysqlDB                  *gorm.DB
	RedisClient              *redis.Client
	KafkaAsyncEventPublisher *kafkautil.EventAsyncPublisher
	KafkaSyncEventPublisher  *kafkautil.EventSyncPublisher
	OutboxDispatcher         *outboxUtil.Dispatcher
	KafkaConsumers           []*kafkautil.Consumer
	BinlogReader             *canalUtil.BinlogReader
	OTELShutdown             OTELShutdown
}

func BuildDependencies(
	httpServer *ginutils.Server,
	mysqlDB *gorm.DB,
	redisClient *redis.Client,
	kafkaAsyncPublisher *kafkautil.EventAsyncPublisher,
	kafkaSyncPublisher *kafkautil.EventSyncPublisher,
	outboxDispatcher *outboxUtil.Dispatcher,
	consumers []*kafkautil.Consumer,
	binlogReader *canalUtil.BinlogReader,
	OTELShutdown OTELShutdown,
) (*Dependencies, error) {
	deps := &Dependencies{
		HttpServer:               httpServer,
		MysqlDB:                  mysqlDB,
		RedisClient:              redisClient,
		KafkaAsyncEventPublisher: kafkaAsyncPublisher,
		KafkaSyncEventPublisher:  kafkaSyncPublisher,
		OutboxDispatcher:         outboxDispatcher,
		KafkaConsumers:           consumers,
		BinlogReader:             binlogReader,
		OTELShutdown:             OTELShutdown,
	}
	return deps, nil
}
