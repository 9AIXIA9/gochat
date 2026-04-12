package di

import (
	canalUtil "gochat/internal/infrastructure/canal"
	ginutils "gochat/internal/infrastructure/gin"
	kafkautil "gochat/internal/infrastructure/kafka"
	gomailUtil "gochat/internal/notification/infrastructure/gomail"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	HttpServer          *ginutils.Server
	MysqlDB             *gorm.DB
	RedisClient         *redis.Client
	KafkaEventPublisher *kafkautil.EventPublisher
	KafkaConsumers      []*kafkautil.Consumer
	BinlogReader        *canalUtil.BinlogReader
	EmailNotifier       *gomailUtil.EmailNotifier
	OTELShutdown        OTELShutdown
}

func BuildDependencies(
	httpServer *ginutils.Server,
	mysqlDB *gorm.DB,
	redisClient *redis.Client,
	kafkaPublisher *kafkautil.EventPublisher,
	consumers []*kafkautil.Consumer,
	binlogReader *canalUtil.BinlogReader,
	emailNotifier *gomailUtil.EmailNotifier,
	OTELShutdown OTELShutdown,
	_ emailServiceAvailable,
) (*Dependencies, error) {
	deps := &Dependencies{
		HttpServer:          httpServer,
		MysqlDB:             mysqlDB,
		RedisClient:         redisClient,
		KafkaEventPublisher: kafkaPublisher,
		KafkaConsumers:      consumers,
		BinlogReader:        binlogReader,
		EmailNotifier:       emailNotifier,
		OTELShutdown:        OTELShutdown,
	}
	return deps, nil
}
