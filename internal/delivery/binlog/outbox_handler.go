package binlog

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
)

//TODO 为什么一执行就是从头开始读，而不是从最新的位点开始读？

var _ canal.EventHandler = (*outboxHandler)(nil)

type outboxHandler struct {
	canal.DummyEventHandler
}

func NewOutboxHandler() canal.EventHandler {
	return &outboxHandler{
		DummyEventHandler: canal.DummyEventHandler{},
	}
}

func (h *outboxHandler) OnRow(e *canal.RowsEvent) error {
	// only care insert into unpublished_events
	if e.Action != canal.InsertAction || e.Table == nil || e.Table.Name != "unpublished_events" {
		return nil
	}

	//TODO 调用用例执行逻辑
	fmt.Println("New unpublished event inserted:", e.Rows)

	return nil
}
