package websocket

import (
	"context"
	"fmt"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
)

//TODO 完善

var _ application.MessageNotifier = (*Manager)(nil)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Notify(_ context.Context, message *domain.Message) error {
	fmt.Println("websocket notify user:", message.Recipient(), "message:", message.Content())
	return nil
}
