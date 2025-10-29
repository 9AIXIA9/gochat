package canal

import (
	"context"
	"strings"

	"gochat/internal/shared/event"

	"github.com/go-mysql-org/go-mysql/canal"
	"go.uber.org/zap"
)

// OutboxConsumer watches MySQL binlog and publishes rows inserted into gochat.unpublished_events
// using the provided event.Publisher. This is an infrastructure adapter implementing CDC for the outbox table.
type OutboxConsumer struct {
	publisher event.Publisher
	lister    event.UnpublishedLister
	canal     *canal.Canal
}

func NewOutboxConsumer(config *BinlogReaderConfig, publisher event.Publisher, lister event.UnpublishedLister) (*OutboxConsumer, error) {
	canalConfig := config.ToCanal()

	cn, err := canal.NewCanal(canalConfig)
	if err != nil {
		return nil, err
	}

	return &OutboxConsumer{
		publisher: publisher,
		lister:    lister,
		canal:     cn,
	}, nil
}

func (c *OutboxConsumer) Start(ctx context.Context) {
	c.canal.SetEventHandler(&outboxHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		publisher:         c.publisher,
		lister:            c.lister,
	})

	// start from latest master position
	go func() {
		if err := c.canal.Run(); err != nil && !strings.Contains(err.Error(), "context canceled") {
			zap.L().Error("binlog canal run failed", zap.Error(err))
		}
	}()
	go func() {
		<-ctx.Done()
		c.canal.Close()
	}()
}

type outboxHandler struct {
	canal.DummyEventHandler
	publisher event.Publisher
	lister    event.UnpublishedLister
}

func (h *outboxHandler) OnRow(e *canal.RowsEvent) error {
	// only care insert into unpublished_events
	if e.Action != canal.InsertAction || e.Table == nil || e.Table.Name != "unpublished_events" {
		return nil
	}

	events, err := h.lister.UnpublishedList(context.Background())
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	return h.publisher.PublishEvents(events)
}
