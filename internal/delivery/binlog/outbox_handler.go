package binlog

import (
	"context"
	"gochat/internal/application"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
)

const timeout = 5 * time.Second

var _ canal.EventHandler = (*outboxHandler)(nil)

type outboxHandler struct {
	canal.DummyEventHandler
	uc application.UnpublishedEventsCreatedUseCase
}

func NewOutboxHandler(uc application.UnpublishedEventsCreatedUseCase) canal.EventHandler {
	return &outboxHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		uc:                uc,
	}
}

func (h *outboxHandler) OnRow(e *canal.RowsEvent) error {
	// only care insert into unpublished_events
	if e.Action != canal.InsertAction || e.Table == nil || e.Table.Name != "unpublished_events" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := h.uc.Execute(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
