package gomail

import (
	"context"
	"gochat/internal/notification/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

//TODO: 除去异步发送，直接由kafka调用（错误信息如果包含 login failed 转换为 服务不可用）
// 提供一个熔断机制，防止邮件服务不可用时大量请求堆积

var _ application.WelcomeEmailNotifier = (*EmailNotifier)(nil)

const (
	maxWaitTime  = 300 * time.Second
	maxMailCache = 100
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

func NewEmailNotifier(senderName string, dialer *gomail.Dialer) *EmailNotifier {
	return &EmailNotifier{
		dialer: dialer,
		taskGenerator: &taskGenerator{
			senderName: senderName,
			address:    dialer.Username,
		},
		taskChan: make(chan *task, maxMailCache),
		closed:   false,
	}
}

func (n *EmailNotifier) NotifyWelcomeEmail(ctx context.Context, email kernel.Email, number kernel.UserNumber) error {
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
		go func() {
			for t := range n.taskChan {
				if err := n.dialer.DialAndSend(t.message); err != nil {
					zap.L().Error(
						"WelcomeEmailNotifier send error",
						zap.String("email", t.email.String()),
						zap.Error(err),
					)
				}
			}
		}()
	})
}

func (n *EmailNotifier) Close() {
	n.closeOnce.Do(func() {
		n.mu.Lock()
		n.closed = true
		close(n.taskChan)
		n.mu.Unlock()
	})
}
