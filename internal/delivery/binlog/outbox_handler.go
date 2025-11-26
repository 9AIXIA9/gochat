package binlog

import (
	"context"
	"fmt"
	"gochat/internal/application"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
)

//TODO 为什么一执行就是从头开始读，而不是从最新的位点开始读？

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

	fmt.Println("OutboxHandler: Detected new unpublished_events insertion", e.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := h.uc.Execute(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
