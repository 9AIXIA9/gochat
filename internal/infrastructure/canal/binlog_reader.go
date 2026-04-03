package canal

import (
	"context"
	"gochat/internal/infrastructure/metrics"
	"gochat/pkg/concurrency"
	"strings"

	"github.com/go-mysql-org/go-mysql/canal"
	"go.uber.org/zap"
)

type BinlogReader struct {
	canal *canal.Canal
}

func NewBinlogReader(
	c *canal.Canal,
	handler canal.EventHandler,
) *BinlogReader {
	reader := &BinlogReader{
		canal: c,
	}
	reader.canal.SetEventHandler(handler)
	return reader
}

func (b *BinlogReader) Start() {
	//从最新的主位置开始
	concurrency.GoSafe(func() {
		pos, err := b.canal.GetMasterPos()
		if err != nil {
			metrics.BinlogReaderRun(context.Background(), "failed", "get_master_pos")
			zap.L().Error("get master binlog position failed", zap.Error(err))
			return
		}
		metrics.BinlogReaderRun(context.Background(), "started", "")
		if err := b.canal.RunFrom(pos); err != nil && !strings.Contains(err.Error(), "context canceled") {
			metrics.BinlogReaderRun(context.Background(), "failed", "run_from")
			zap.L().Error("binlog canal run failed", zap.Error(err))
		}
	})
}

func (b *BinlogReader) Close() {
	b.canal.Close()
}
