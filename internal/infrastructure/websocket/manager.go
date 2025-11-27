package websocket

import (
	myErrors "gochat/internal/shared/errors"
	"sync"

	"gochat/internal/shared/kernel"
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
	m.mu.Lock()
	defer m.mu.Unlock()
	// 如果已存在旧连接，先关掉
	if old, ok := m.clients[id]; ok {
		old.Close()
	}
	m.clients[id] = c
}

func (m *Manager) Unregister(id kernel.UserID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.clients[id]; ok {
		c.Close()
		delete(m.clients, id)
	}
}

func (m *Manager) SendTo(id kernel.UserID, resp *Response) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.clients[id]; ok {
		if err := c.SendResponse(resp); err != nil {
			return err
		}
		return nil
	}
	return myErrors.ErrNotFound
}

func (m *Manager) Broadcast(ids []kernel.UserID, resp *Response) []kernel.UserID {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := 0

	idsSuccess := make([]kernel.UserID, 0, len(ids))
	for _, id := range ids {
		if c, ok := m.clients[id]; ok {
			if err := c.SendResponse(resp); err != nil {
				continue
			}
			idsSuccess = append(idsSuccess, id)
			n++
		}
	}

	return idsSuccess
}
