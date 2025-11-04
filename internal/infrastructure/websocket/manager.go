package websocket

import (
	"context"
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

func (m *Manager) Enqueue(_ context.Context, _ kernel.UserID, _ *domain.Message) error {
	return nil
}
