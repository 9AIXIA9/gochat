package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/metrics"
	myErrors "gochat/internal/shared/errors"
	"sync"

	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

const defaultBaseClientsCount = 10000

type Manager struct {
	mu      sync.RWMutex
	clients map[kernel.UserID]*Client
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[kernel.UserID]*Client, defaultBaseClientsCount),
	}
}

func (m *Manager) Register(id kernel.UserID, c *Client) {
	// Close outside of the lock to avoid holding the mutex while closing.
	var toClose *Client

	m.mu.Lock()
	if old, ok := m.clients[id]; ok {
		toClose = old
	}
	m.clients[id] = c
	m.mu.Unlock()

	if toClose != nil {
		metrics.WSDisconnect(context.Background(), "replaced")
		toClose.Close()
	}
	metrics.WSConnectionDelta(context.Background(), 1)
}

func (m *Manager) Unregister(id kernel.UserID) {
	m.mu.Lock()
	c, ok := m.clients[id]
	if ok {
		delete(m.clients, id)
	}
	m.mu.Unlock()

	if ok {
		metrics.WSConnectionDelta(context.Background(), -1)
		metrics.WSDisconnect(context.Background(), "unregister")
		zap.L().Debug(
			"websocket manager: unregistered client",
			zap.String("userID", id.String()),
		)
		c.Close()
	}
}

// UnregisterClient removes a client by pointer and returns its user ID if found.
func (m *Manager) UnregisterClient(target *Client) (id kernel.UserID, found bool) {
	m.mu.Lock()
	for uid, c := range m.clients {
		if c == target {
			delete(m.clients, uid)
			id = uid
			found = true
			break
		}
	}
	m.mu.Unlock()

	if found {
		metrics.WSConnectionDelta(context.Background(), -1)
		metrics.WSDisconnect(context.Background(), "unregister_by_pointer")
		zap.L().Debug("websocket manager: unregistered client by pointer", zap.String("userID", id.String()))
	}
	return
}

func (m *Manager) SendTo(id kernel.UserID, msg *Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		metrics.WSMessageOut(context.Background(), "manager_sendto", "marshal_failed")
		return err
	}

	m.mu.RLock()
	c, ok := m.clients[id]
	m.mu.RUnlock()
	if ok {
		if err := c.Send(msgBytes); err != nil {
			metrics.WSMessageOut(context.Background(), "manager_sendto", "send_failed")
			return err
		}
		metrics.WSMessageOut(context.Background(), "manager_sendto", "ok")
		return nil
	}
	metrics.WSMessageOut(context.Background(), "manager_sendto", "not_found")
	return myErrors.ErrNotFound
}

func (m *Manager) Broadcast(ids []kernel.UserID, msg *Message) ([]kernel.UserID, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		metrics.WSMessageOut(context.Background(), "manager_broadcast", "marshal_failed")
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	idsSuccess := make([]kernel.UserID, 0, len(ids))
	for _, id := range ids {
		if c, ok := m.clients[id]; ok {
			if err := c.Send(msgBytes); err != nil {
				metrics.WSMessageOut(context.Background(), "manager_broadcast", "send_failed")
				continue
			}
			metrics.WSMessageOut(context.Background(), "manager_broadcast", "ok")
			idsSuccess = append(idsSuccess, id)
		} else {
			metrics.WSMessageOut(context.Background(), "manager_broadcast", "not_found")
		}
	}

	return idsSuccess, nil
}
