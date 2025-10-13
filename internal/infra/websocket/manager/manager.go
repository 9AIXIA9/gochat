package manager

import (
	"github.com/gorilla/websocket"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/client"
	"gochat/internal/types"

	"sync"
)

var _ domain.MessageSender = NewWebsocketManager()

const (
	clientsCache = 200
)

type Manager struct {
	clients map[domain.UserNumber]client.Client
	mu      sync.RWMutex
}

func NewWebsocketManager() *Manager {
	return &Manager{
		clients: make(map[domain.UserNumber]client.Client, clientsCache),
	}
}

func (m *Manager) AddClient(conn *websocket.Conn, number domain.UserNumber) func() {
	c := client.New(conn)

	// 先注册到 map（不在持锁状态下启动协程）
	m.mu.Lock()
	old, ok := m.clients[number]
	m.clients[number] = c
	m.mu.Unlock()

	// 如有旧连接，移除并关闭（避免在持锁时做关闭）
	if ok && old != nil {
		old.Close()
	}

	c.Start()
	return func() {
		c.Wait()

		m.mu.Lock()
		delete(m.clients, number)
		m.mu.Unlock()

		c.Close()
	}
}

func (m *Manager) SendMessage(number domain.UserNumber, msg *domain.Message) error {
	m.mu.RLock()
	c, ok := m.clients[number]
	m.mu.RUnlock()

	if ok {
		return c.Send(msg)
	}
	return types.ErrNotFound
}
