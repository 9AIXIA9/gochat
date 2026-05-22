package outbox

import (
	"context"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var _ Triggerer = (*Dispatcher)(nil)

type executor interface {
	Execute(context.Context, *kernel.NoInput) (*kernel.NoOutput, error)
}

// Triggerer receives a non-blocking signal to start an outbox drain.
type Triggerer interface {
	Trigger()
}

type Dispatcher struct {
	executor      executor
	sweepInterval time.Duration
	triggerCh     chan struct{}
	stopCh        chan struct{}

	started atomic.Bool
	running atomic.Bool

	stopOnce sync.Once
	wg       sync.WaitGroup
}

func NewDispatcher(executor executor, sweepInterval time.Duration, triggerBuffer int) (*Dispatcher, error) {
	if err := validate.NotNil(executor); err != nil {
		return nil, err
	}
	if sweepInterval <= 0 {
		sweepInterval = time.Second
	}
	if triggerBuffer <= 0 {
		triggerBuffer = 1
	}

	return &Dispatcher{
		executor:      executor,
		sweepInterval: sweepInterval,
		triggerCh:     make(chan struct{}, triggerBuffer),
		stopCh:        make(chan struct{}),
	}, nil
}

func (d *Dispatcher) Start(ctx context.Context) {
	if d == nil || !d.started.CompareAndSwap(false, true) {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		defer zap.L().Info("outbox dispatcher stopped")

		ticker := time.NewTicker(d.sweepInterval)
		defer ticker.Stop()

		d.Trigger()

		for {
			select {
			case <-ctx.Done():
				return
			case <-d.stopCh:
				return
			case <-ticker.C:
				d.runOnce(ctx, "ticker")
			case <-d.triggerCh:
				d.runOnce(ctx, "trigger")
			}
		}
	}()
}

func (d *Dispatcher) Trigger() {
	if d == nil {
		return
	}
	select {
	case d.triggerCh <- struct{}{}:
	default:
	}
}

func (d *Dispatcher) Stop() {
	if d == nil {
		return
	}
	d.stopOnce.Do(func() {
		close(d.stopCh)
	})
	d.wg.Wait()
}

func (d *Dispatcher) runOnce(ctx context.Context, source string) {
	if d.running.Load() {
		return
	}
	if !d.running.CompareAndSwap(false, true) {
		return
	}
	defer d.running.Store(false)

	start := time.Now()
	_, err := d.executor.Execute(ctx, nil)
	if err != nil {
		zap.L().Warn("outbox dispatcher drain failed",
			zap.String("source", source),
			zap.Duration("elapsed", time.Since(start)),
			zap.Error(err),
		)
		return
	}

	zap.L().Info("outbox dispatcher drain completed",
		zap.String("source", source),
		zap.Duration("elapsed", time.Since(start)),
	)
}
