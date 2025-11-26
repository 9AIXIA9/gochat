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
) *BinlogReader {
	return &BinlogReader{
		canal: c,
	}
}

func (b *BinlogReader) SetHandler(handler canal.EventHandler) {
	b.canal.SetEventHandler(handler)
}

func (b *BinlogReader) Start() {
	// start from latest master position
	utils.GoSafe(func() {
		if err := b.canal.Run(); err != nil && !strings.Contains(err.Error(), "context canceled") {
			zap.L().Error("binlog canal run failed", zap.Error(err))
		}
	})
}

func (b *BinlogReader) Close() {
	b.canal.Close()
}
