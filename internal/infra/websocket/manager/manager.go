package manager

import (
	"errors"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/client"
	"gochat/internal/types"
	"sync"
)

const (
	clientsCache = 200
)

type Manager struct {
	clients map[domain.UserNumber]client.Client
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[domain.UserNumber]client.Client, clientsCache),
	}
}

func (m *Manager) AddClient(c client.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.Number()] = c

	c.SetCloseHandler(func(code int, text string) error {
		m.DropClient(c.Number())
		return nil
	})
}

func (m *Manager) DropClient(number domain.UserNumber) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, number)
}

func (m *Manager) Send(number domain.UserNumber, msg *domain.Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[number]
	if ok {
		return c.Write(msg)
	}
	return types.ErrNotFound
}

func (m *Manager) SendUserManyMsgs(number domain.UserNumber, msgs []*domain.Message) []domain.MessageID {
	msgIDs := make([]domain.MessageID, 0, len(msgs))
	for _, msg := range msgs {
		err := m.Send(number, msg)
		if err == nil {
			msgIDs = append(msgIDs, msg.ID())
		}
		if err != nil && !errors.Is(err, types.ErrNotFound) {
			zap.L().Error("send message failed",
				zap.Int64("user_number", int64(number)),
				zap.String("message_id", string(msg.ID())),
				zap.Error(err))
		}
	}
	return msgIDs
}

func (m *Manager) SendMsgToManyUsers(msg *domain.Message, numbers []domain.UserNumber) []domain.UserNumber {
	numbersSent := make([]domain.UserNumber, 0, len(numbers))
	for _, number := range numbers {
		err := m.Send(number, msg)
		if err == nil {
			numbersSent = append(numbersSent, number)
		}
		if err != nil && !errors.Is(err, types.ErrNotFound) {
			zap.L().Error("send message failed",
				zap.Int64("user_number", int64(number)),
				zap.String("message_id", string(msg.ID())),
				zap.Error(err))
		}
	}
	return numbersSent
}
