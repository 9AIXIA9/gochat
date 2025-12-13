package websocket

import "context"

type Topic string

func (t Topic) String() string {
	return string(t)
}

const websocketTopicKey = "topic"

func GetTopic(ctx context.Context) Topic {
	topic, ok := ctx.Value(websocketTopicKey).(Topic)
	if !ok {
		return ""
	}
	return topic
}

func SetTopic(ctx context.Context, topic Topic) context.Context {
	return context.WithValue(ctx, websocketTopicKey, topic)
}
