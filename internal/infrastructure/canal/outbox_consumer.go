package canal

import (
	"context"
	"time"

	"gochat/internal/infrastructure/prometheus"
	"gochat/pkg/utils"
	"strings"

	"gochat/internal/shared/event"

	"github.com/go-mysql-org/go-mysql/canal"
	"go.uber.org/zap"
)

// OutboxConsumer watches MySQL binlog and publishes rows inserted into gochat.unpublished_events
// using the provided event.Publisher. This is an infrastructure adapter implementing CDC for the outbox table.
type OutboxConsumer struct {
	publisher event.Publisher
	lister    event.UnpublishedEventsLister
	canal     *canal.Canal
	metrics   *prometheus.Metrics
}

func NewOutboxConsumer(c *canal.Canal, publisher event.Publisher, lister event.UnpublishedEventsLister) *OutboxConsumer {
	return &OutboxConsumer{
		publisher: publisher,
		lister:    lister,
		canal:     c,
	}
}

func (c *OutboxConsumer) SetMetrics(m *prometheus.Metrics) { c.metrics = m }

func (c *OutboxConsumer) Start() {
	c.canal.SetEventHandler(&outboxHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		publisher:         c.publisher,
		lister:            c.lister,
		metrics:           c.metrics,
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
	publisher event.Publisher
	lister    event.UnpublishedEventsLister
	metrics   *prometheus.Metrics
}

func (h *outboxHandler) OnRow(e *canal.RowsEvent) error {
	// only care insert into unpublished_events
	if e.Action != canal.InsertAction || e.Table == nil || e.Table.Name != "unpublished_events" {
		return nil
	}
	if h.metrics != nil {
		h.metrics.OutboxPolled.Inc()
	}

	events, err := h.lister.ListUnpublishedEvents(context.Background())
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}
	if h.metrics != nil {
		h.metrics.OutboxBatchSize.Observe(float64(len(events)))
	}

	start := time.Now()
	if err := h.publisher.Publish(events); err != nil {
		if h.metrics != nil {
			h.metrics.OutboxPublishes.WithLabelValues("error").Inc()
			h.metrics.OutboxDur.WithLabelValues("error").Observe(time.Since(start).Seconds())
		}
		return err
	}
	if h.metrics != nil {
		h.metrics.OutboxPublishes.WithLabelValues("success").Inc()
		h.metrics.OutboxDur.WithLabelValues("success").Observe(time.Since(start).Seconds())
	}
	return nil
}
