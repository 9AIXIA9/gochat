package di

import (
	"gochat/config"
	"gochat/internal/delivery/binlog"
	canalInfra "gochat/internal/infrastructure/canal"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/google/wire"
)

type BinlogReaderHandlerEnsured bool

var BinlogSet = wire.NewSet(
	provideCanal,
	provideCanalBinlogReader,
	provideCanalBinlogReaderHandlerEnsured,
)

func provideCanal(appConfig *config.App) (*canal.Canal, error) {
	return canalInfra.NewCanal(appConfig.BinlogReader)
}

func provideCanalBinlogReader(c *canal.Canal) *canalInfra.BinlogReader {
	return canalInfra.NewBinlogReader(c)
}

func provideCanalBinlogReaderHandlerEnsured(reader *canalInfra.BinlogReader) BinlogReaderHandlerEnsured {
	reader.SetHandler(binlog.NewOutboxHandler())
	return true
}
