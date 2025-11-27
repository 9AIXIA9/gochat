package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"

	"go.uber.org/zap"
)

func NewNotFoundHandler() websocket.HandlerFunc {
	return func(_ context.Context, data []byte) ([]byte, error) {
		zap.L().Debug("websocket: not found handler invoked", zap.ByteString("data", data))
		return json.Marshal(&websocket.MessageData{
			Message: "topic is not found",
		})
	}
}
