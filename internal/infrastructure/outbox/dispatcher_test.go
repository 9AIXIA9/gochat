package outbox

import (
	"context"
	"gochat/internal/shared/kernel"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeExecutor struct {
	calls   int32
	enter   chan struct{}
	release chan struct{}
	once    sync.Once
}

func (f *fakeExecutor) Execute(ctx context.Context, _ *kernel.NoInput) (*kernel.NoOutput, error) {
	atomic.AddInt32(&f.calls, 1)
	if f.enter != nil {
		f.once.Do(func() { close(f.enter) })
	}
	if f.release != nil {
		select {
		case <-f.release:
		case <-ctx.Done():
		}
	}
	return nil, nil
}

func TestDispatcherTriggerRunsExecutor(t *testing.T) {
	exec := &fakeExecutor{enter: make(chan struct{})}
	d, err := NewDispatcher(exec, 10*time.Millisecond, 1)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.Start(ctx)
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&exec.calls) > 0
	}, time.Second, 10*time.Millisecond)

	d.Stop()
	<-exec.enter
}

func TestDispatcherSingleRunGuard(t *testing.T) {
	release := make(chan struct{})
	exec := &fakeExecutor{enter: make(chan struct{}), release: release}
	d, err := NewDispatcher(exec, 10*time.Millisecond, 1)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.Start(ctx)
	<-exec.enter

	d.Trigger()
	d.Trigger()
	time.Sleep(30 * time.Millisecond)
	require.Equal(t, int32(1), atomic.LoadInt32(&exec.calls))

	cancel()
	close(release)
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&exec.calls) >= 1
	}, time.Second, 10*time.Millisecond)

	d.Stop()
}
