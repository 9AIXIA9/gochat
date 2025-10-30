package gomail

import (
	"context"
	"fmt"
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

var _ application.EmailNotifier = (*EmailNotifier)(nil)

const (
	maxWorkers         = 3
	maxRetries         = 3
	maxBackoffDuration = 30 * time.Second
	maxWaitTime        = 300 * time.Second
	maxNoticeCache     = 100
)

type EmailNotifier struct {
	senderName string

	dialer *gomail.Dialer

	taskChan  chan *task
	mu        sync.RWMutex
	startOnce sync.Once
	closeOnce sync.Once
	closed    bool
}

func NewEmailNotifier(senderName string, config *EmailNotifierConfig) *EmailNotifier {
	return &EmailNotifier{
		senderName: senderName,
		dialer:     gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
		taskChan:   make(chan *task, maxNoticeCache),
		mu:         sync.RWMutex{},
		closeOnce:  sync.Once{},
		closed:     false,
	}
}

func (n *EmailNotifier) Enqueue(ctx context.Context, email kernel.Email, notice *domain.Notice, onSuccess func() error) error {
	fmt.Println("EmailNotifier Enqueue called")
	n.mu.RLock()
	if n.closed {
		n.mu.RUnlock()
		return myErrors.ErrHasBeenClosed
	}

	ctx, cancel := context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	select {
	case n.taskChan <- &task{notice: notice, email: email, onSuccess: onSuccess}:
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
			utils.GoSafe(n.startWorker)
		}
	})
}

func (n *EmailNotifier) startWorker() {
	for t := range n.taskChan {
		if err := n.sendWithRetry(t); err == nil && t.onSuccess != nil {
			if err := t.onSuccess(); err != nil {
				zap.L().Error(
					"EmailNotifier onSuccess callback error",
					zap.String("notice_id", string(t.notice.ID())),
					zap.Error(err),
				)
			}
		}
	}
}

func (n *EmailNotifier) sendWithRetry(t *task) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			utils.BackoffWait(maxBackoffDuration, attempt)
		}

		msg, buildErr := n.composeMessage(t)
		if buildErr != nil {
			// 构造错误为非重试错误
			return buildErr
		}

		if err := n.dialer.DialAndSend(msg); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	return lastErr
}

func (n *EmailNotifier) composeMessage(t *task) (*gomail.Message, error) {
	if t == nil {
		return nil, myErrors.ErrEmptyPointer
	}

	m := gomail.NewMessage()

	// From
	fromAddr := n.dialer.Username
	if n.senderName != "" {
		m.SetAddressHeader("From", fromAddr, n.senderName)
	} else {
		m.SetHeader("From", fromAddr)
	}

	// To
	m.SetHeader("To", fmt.Sprint(t.email))

	// Subject
	subject := fmt.Sprintf("Notice %s", fmt.Sprint(t.notice.Theme()))
	m.SetHeader("Subject", subject)

	// Body
	m.SetBody("text/html", t.notice.Content())

	return m, nil
}

func (n *EmailNotifier) Close() {
	n.closeOnce.Do(func() {
		n.mu.Lock()
		n.closed = true
		close(n.taskChan)
		n.mu.Unlock()
	})
}
