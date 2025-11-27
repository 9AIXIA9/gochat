package utils

import (
	"context"
	"gochat/internal/shared/kernel"
)

const (
	userIDKey         = "user_id"
	websocketTopicKey = "topic"
)

func SetUserID(ctx context.Context, userID kernel.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func SetWebsocketTopic(ctx context.Context, topic string) context.Context {
	return context.WithValue(ctx, websocketTopicKey, topic)
}

func GetUserID(ctx context.Context) kernel.UserID {
	userID, ok := ctx.Value(userIDKey).(kernel.UserID)
	if !ok {
		return ""
	}
	return userID
}

func GetWebsocketTopic(ctx context.Context) string {
	topic, ok := ctx.Value(websocketTopicKey).(string)
	if !ok {
		return ""
	}
	return topic
}
