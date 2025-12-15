package websocket

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"sync"

	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

type Manager struct {
	mu      sync.RWMutex
	clients map[kernel.UserID]*Client
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[kernel.UserID]*Client),
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
		toClose.Close()
	}
}

func (m *Manager) Unregister(id kernel.UserID) {
	m.mu.Lock()
	c, ok := m.clients[id]
	if ok {
		delete(m.clients, id)
	}
	m.mu.Unlock()

	if ok {
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
		zap.L().Debug("websocket manager: unregistered client by pointer", zap.String("userID", id.String()))
	}
	return
}

func (m *Manager) SendTo(id kernel.UserID, resp *Response) error {
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	m.mu.RLock()
	c, ok := m.clients[id]
	m.mu.RUnlock()
	if ok {
		if err := c.Send(respData); err != nil {
			return err
		}
		return nil
	}
	return myErrors.ErrNotFound
}

func (m *Manager) Broadcast(ids []kernel.UserID, resp *Response) ([]kernel.UserID, error) {
	respData, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	idsSuccess := make([]kernel.UserID, 0, len(ids))
	for _, id := range ids {
		if c, ok := m.clients[id]; ok {
			if err := c.Send(respData); err != nil {
				continue
			}
			idsSuccess = append(idsSuccess, id)
		}
	}

	return idsSuccess, nil
}
