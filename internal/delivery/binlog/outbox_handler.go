package binlog

import (
	"context"
	"gochat/internal/infrastructure/metrics"
	"gochat/internal/infrastructure/outbox"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
)

var _ canal.EventHandler = (*outboxHandler)(nil)

type outboxHandler struct {
	canal.DummyEventHandler
	triggerer outbox.Triggerer
}

func NewOutboxHandler(triggerer outbox.Triggerer) canal.EventHandler {
	return &outboxHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		triggerer:         triggerer,
	}
}

func (h *outboxHandler) OnRow(e *canal.RowsEvent) error {
	// only care insert into unpublished_events
	if e.Action != canal.InsertAction || e.Table == nil || e.Table.Name != "unpublished_events" {
		return nil
	}

	start := time.Now().UTC()

	if h.triggerer != nil {
		h.triggerer.Trigger()
	}

	ctx := context.Background()
	metrics.OutboxProcess(ctx, "binlog_trigger", "ok", "signal_sent")
	metrics.OutboxProcessDuration(ctx, "binlog_trigger", "ok", time.Since(start).Seconds())

	return nil
}
