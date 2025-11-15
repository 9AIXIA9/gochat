package canal

import (
	"context"
	"gochat/pkg/utils"
	"strings"

	"gochat/internal/shared/event"

	"github.com/go-mysql-org/go-mysql/canal"
	"go.uber.org/zap"
)

// OutboxConsumer watches MySQL binlog and publishes rows inserted into gochat.unpublished_events
// using the provided event.ManyPublisher. This is an infrastructure adapter implementing CDC for the outbox table.
type OutboxConsumer struct {
	publisher event.ManyPublisher
	lister    event.UnpublishedLister
	canal     *canal.Canal
}

func NewOutboxConsumer(c *canal.Canal, publisher event.ManyPublisher, lister event.UnpublishedLister) *OutboxConsumer {
	return &OutboxConsumer{
		publisher: publisher,
		lister:    lister,
		canal:     c,
	}
}

func (c *OutboxConsumer) Start() {
	c.canal.SetEventHandler(&outboxHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		publisher:         c.publisher,
		lister:            c.lister,
	})

	// start from latest master position
	utils.GoSafe(func() {
		if err := c.canal.Run(); err != nil && !strings.Contains(err.Error(), "context canceled") {
			zap.L().Error("binlog canal run failed", zap.Error(err))
		}
	})
}

func (c *OutboxConsumer) Close() {
	c.canal.Close()
}

type outboxHandler struct {
	canal.DummyEventHandler
	publisher event.ManyPublisher
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

	return h.publisher.Publishes(events)
}
