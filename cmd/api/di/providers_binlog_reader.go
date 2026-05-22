package di

import (
	"gochat/config"
	"gochat/internal/delivery/binlog"
	canalInfra "gochat/internal/infrastructure/canal"
	outboxInfra "gochat/internal/infrastructure/outbox"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/google/wire"
)

var BinlogSet = wire.NewSet(
	provideCanal,
	provideCanalBinlogReader,
	provideCanalBinlogReaderHandler,
)

func provideCanal(appConfig *config.App) (*canal.Canal, error) {
	return canalInfra.NewCanal(appConfig.BinlogReader)
}

func provideCanalBinlogReader(c *canal.Canal, handler canal.EventHandler) *canalInfra.BinlogReader {
	return canalInfra.NewBinlogReader(c, handler)
}

func provideCanalBinlogReaderHandler(dispatcher *outboxInfra.Dispatcher) canal.EventHandler {
	return binlog.NewOutboxHandler(dispatcher)
}
