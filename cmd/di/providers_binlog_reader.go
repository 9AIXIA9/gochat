package di

import (
	"gochat/config"
	"gochat/internal/application"
	"gochat/internal/delivery/binlog"
	canalInfra "gochat/internal/infrastructure/canal"

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

func provideCanalBinlogReaderHandler(uc application.UnpublishedEventsCreatedUseCase) canal.EventHandler {
	return binlog.NewOutboxHandler(uc)
}
