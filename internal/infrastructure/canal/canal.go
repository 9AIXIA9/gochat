package canal

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
)

func NewCanal(config *BinlogReaderConfig) (*canal.Canal, error) {
	canalConfig := canal.NewDefaultConfig()
	canalConfig.User = config.User
	canalConfig.Addr = config.Addr
	canalConfig.Password = config.Password
	canalConfig.Dump.TableDB = config.TableDB
	canalConfig.Dump.ExecutionPath = "" // 不使用 mysqldump 工具
	canalConfig.Flavor = mysql.MySQLFlavor
	canalConfig.IncludeTableRegex = []string{".*\\.unpublished_events"}
	canalConfig.DiscardNoMetaRowEvent = true
	canalConfig.HeartbeatPeriod = 200 * time.Millisecond
	canalConfig.ReadTimeout = 300 * time.Millisecond
	canalConfig.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	cn, err := canal.NewCanal(canalConfig)
	if err != nil {
		return nil, err
	}
	return cn, nil
}
