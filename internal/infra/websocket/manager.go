package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gochat/internal/domain"
	"gochat/internal/types"
	"net/http"
	"sync"
)

const (
	SendCache         = 256
	PreparedRoomCount = 32
)

type Manager struct {
	websocket.Upgrader
	roomsMutex sync.RWMutex
	rooms      map[domain.RoomNumber]*Room
}

func NewManager() *Manager {
	return &Manager{
		Upgrader:   websocket.Upgrader{},
		roomsMutex: sync.RWMutex{},
		rooms:      make(map[domain.RoomNumber]*Room, PreparedRoomCount),
	}
}

// CreateRoom 创建房间
func (m *Manager) CreateRoom(roomNumber domain.RoomNumber) error {
	//检查是否存在此房间
	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	_, exist := m.rooms[roomNumber]
	if exist {
		return types.ErrDuplicateKey
	}

	//创建房间
	room := &Room{
		broadcastChan: make(chan []byte, SendCache),
		clients:       nil,
		clientsMutex:  sync.RWMutex{},
	}

	m.rooms[roomNumber] = room
	return nil
}

// EstablishConnection 建立连接
func (m *Manager) EstablishConnection(w http.ResponseWriter, r *http.Request, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查是否存在此房间
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()

	room, exist := m.rooms[roomNumber]
	if !exist {
		return types.ErrNotFound
	}

	//检查用户是否已进入
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	_, exist = room.clients[userNumber]
	if exist {
		return types.ErrDuplicateKey
	}

	// 升级HTTP连接到WebSocket
	ws, err := m.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("websocket upgrade error:%w", err)
	}

	//创建客户端
	client := Client{
		conn:      ws,
		sendChan:  make(chan []byte, SendCache),
		broadcast: room.broadcastChan,
	}

	room.clients[userNumber] = client

	//持续接收消息
	go client.readPump()
	go client.writePump()
	return nil
}

// ExitConnection 关闭连接
func (m *Manager) ExitConnection(userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查房间是否存在
	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	room, exist := m.rooms[roomNumber]
	if !exist {
		return nil
	}

	//关闭用户连接
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	client, exist := room.clients[userNumber]
	if exist {
		if err := client.close(); err != nil {
			return err
		}

		delete(room.clients, userNumber)

		if len(room.clients) == 0 {
			room.close()
			delete(m.rooms, roomNumber)
		}
	}
	return nil
}
