package domain

import (
	"context"
	"time"
)

// EventTopic 事件主题类型
type EventTopic string

func (e EventTopic) String() string {
	return string(e)
}

type EventID string

func (e EventID) String() string {
	return string(e)
}

// Event 领域事件接口
type Event interface {
	// EventID 获取事件ID
	EventID() EventID

	// EventTopic 获取事件主题
	EventTopic() EventTopic

	// OccurredOn 获取事件发生时间
	OccurredOn() time.Time

	// ToPayload 序列化事件数据
	ToPayload() ([]byte, error)
}

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event Event) error
}

// EventPublisher 事件发布者接口
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

// EventSubscriber 事件订阅者接口
type EventSubscriber interface {
	Subscribe(ctx context.Context, topic EventTopic, handler EventHandler) error
	Unsubscribe(ctx context.Context, topic EventTopic, handlerName string) error
}
