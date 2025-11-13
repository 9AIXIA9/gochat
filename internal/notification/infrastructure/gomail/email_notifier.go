package gomail

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var _ application.WelcomeEmailNotifier = (*EmailNotifier)(nil)

const (
	maxWorkers         = 3
	maxRetries         = 3
	maxBackoffDuration = 30 * time.Second
	maxWaitTime        = 300 * time.Second
	maxMailCache       = 100
)

type EmailNotifier struct {
	dialer *gomail.Dialer

	taskGenerator *taskGenerator

	taskChan  chan *task
	mu        sync.RWMutex
	startOnce sync.Once
	closeOnce sync.Once
	closed    bool
}

func NewEmailNotifier(senderName string, config *EmailNotifierConfig) *EmailNotifier {
	return &EmailNotifier{
		dialer: gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
		taskGenerator: &taskGenerator{
			senderName: senderName,
			address:    config.Username,
		},
		taskChan: make(chan *task, maxMailCache),
		closed:   false,
	}
}

func (n *EmailNotifier) NotifyWelcomeEmail(ctx context.Context, email kernel.Email, number domain.UserNumber) error {
	return n.enqueue(ctx, n.taskGenerator.generateUserCreatedTask(email, number))
}

func (n *EmailNotifier) enqueue(ctx context.Context, task *task) error {
	n.mu.RLock()
	if n.closed {
		n.mu.RUnlock()
		return myErrors.ErrHasBeenClosed
	}

	ctx, cancel := context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	select {
	case n.taskChan <- task:
		n.mu.RUnlock()
		return nil
	case <-ctx.Done():
		n.mu.RUnlock()
		return ctx.Err()
	}
}

func (n *EmailNotifier) Start() {
	n.startOnce.Do(func() {
		for i := 0; i < maxWorkers; i++ {
			go n.startWorker()
		}
	})
}

func (n *EmailNotifier) startWorker() {
	for t := range n.taskChan {
		if err := n.sendWithRetry(t.message); err != nil {
			zap.L().Error(
				"WelcomeEmailNotifier send error",
				zap.String("email", t.email.String()),
				zap.Error(err),
			)
		}
	}
}

func (n *EmailNotifier) sendWithRetry(msg *gomail.Message) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			utils.BackoffWait(maxBackoffDuration, attempt)
		}

		if err := n.dialer.DialAndSend(msg); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	return lastErr
}

func (n *EmailNotifier) Close() {
	n.closeOnce.Do(func() {
		n.mu.Lock()
		n.closed = true
		close(n.taskChan)
		n.mu.Unlock()
	})
}
