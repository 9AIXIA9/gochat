package binlog

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/stretchr/testify/require"
)

type fakeTriggerer struct {
	calls int32
}

func (f *fakeTriggerer) Trigger() {
	atomic.AddInt32(&f.calls, 1)
}

func TestOutboxHandlerOnRowTriggersDispatcher(t *testing.T) {
	triggerer := &fakeTriggerer{}
	h := &outboxHandler{triggerer: triggerer}

	err := h.OnRow(&canal.RowsEvent{Action: canal.InsertAction, Table: &schema.Table{Name: "unpublished_events"}})
	require.NoError(t, err)
	require.Equal(t, int32(1), atomic.LoadInt32(&triggerer.calls))
}

func TestOutboxHandlerOnRowIgnoresOtherRows(t *testing.T) {
	triggerer := &fakeTriggerer{}
	h := &outboxHandler{triggerer: triggerer}

	err := h.OnRow(&canal.RowsEvent{Action: canal.UpdateAction, Table: &schema.Table{Name: "unpublished_events"}})
	require.NoError(t, err)
	require.Equal(t, int32(0), atomic.LoadInt32(&triggerer.calls))
}

func TestOutboxHandlerSignalIsFast(t *testing.T) {
	triggerer := &fakeTriggerer{}
	h := &outboxHandler{triggerer: triggerer}

	start := time.Now().UTC()
	err := h.OnRow(&canal.RowsEvent{Action: canal.InsertAction, Table: &schema.Table{Name: "unpublished_events"}})
	require.NoError(t, err)
	require.Less(t, time.Since(start), 50*time.Millisecond)
}
