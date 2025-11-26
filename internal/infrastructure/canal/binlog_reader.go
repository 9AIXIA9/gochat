package canal

import (
	"gochat/pkg/utils"
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
	utils.GoSafe(func() {
		pos, err := b.canal.GetMasterPos()
		if err != nil {
			zap.L().Error("get master binlog position failed", zap.Error(err))
			return
		}
		if err := b.canal.RunFrom(pos); err != nil && !strings.Contains(err.Error(), "context canceled") {
			zap.L().Error("binlog canal run failed", zap.Error(err))
		}
	})
}

func (b *BinlogReader) Close() {
	b.canal.Close()
}
