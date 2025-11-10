package websocket

import (
	"context"
	"fmt"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

//TODO 完善

var _ application.MessageNotifier = (*Manager)(nil)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Notify(_ context.Context, recipient kernel.UserID, message *domain.Message) error {
	fmt.Println("websocket notify user:", recipient, "message:", message.Content())
	return nil
}
