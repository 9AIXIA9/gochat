package websocket

import (
	myErrors "gochat/internal/shared/errors"
	"net/http"
	"sync"

	"gochat/internal/shared/kernel"

	"github.com/gorilla/websocket"
)

type Manager struct {
	upgrader *websocket.Upgrader

	mu          sync.RWMutex
	connections map[kernel.UserID]*Client
}

func NewManager(upgrader *websocket.Upgrader) *Manager {
	return &Manager{
		upgrader:    upgrader,
		connections: make(map[kernel.UserID]*Client),
	}
}

func (m *Manager) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return m.upgrader.Upgrade(w, r, nil)
}

func (m *Manager) Register(id kernel.UserID, c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 如果已存在旧连接，先关掉
	if old, ok := m.connections[id]; ok {
		old.Close()
	}
	m.connections[id] = c
}

func (m *Manager) Unregister(id kernel.UserID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.connections[id]; ok {
		c.Close()
		delete(m.connections, id)
	}
}

func (m *Manager) Get(id kernel.UserID) (*Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.connections[id]
	return c, ok
}

func (m *Manager) SendTo(id kernel.UserID, b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.connections[id]; ok {
		if err := c.Send(b); err != nil {
			return err
		}
		return nil
	}
	return myErrors.ErrNotFound
}

func (m *Manager) Broadcast(b []byte) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := 0
	for _, c := range m.connections {
		_ = c.Send(b)
		n++
	}
	return n
}
