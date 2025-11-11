package websocket

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var _ application.MessageNotifier = (*Manager)(nil)

const (
	clientsCache = 200
)

type Manager struct {
	upgrader *websocket.Upgrader
	clients  map[kernel.UserID]*client
	mu       sync.RWMutex
}

func NewManager(origins []string) *Manager {
	return &Manager{
		upgrader: newUpgrader(origins),
		clients:  make(map[kernel.UserID]*client, clientsCache),
	}
}

func (m *Manager) Connect(userID kernel.UserID, w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Client, error) {
	conn, err := m.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	c := newClient(conn)
	m.mu.Lock()
	m.clients[userID] = c
	m.mu.Unlock()

	return c, nil
}

func (m *Manager) DisConnect(userID kernel.UserID) {
	c := m.clients[userID]
	m.mu.Lock()
	delete(m.clients, userID)
	m.mu.Unlock()

	//close client
	c.Close()
}

func (m *Manager) Notify(recipient kernel.UserID, message *domain.Message) error {
	m.mu.RLock()
	c, ok := m.clients[recipient]
	m.mu.RUnlock()

	if !ok {
		return myErrors.ErrNotFound
	}

	if err := c.Send(message); err != nil {
		return err
	}

	return nil
}
