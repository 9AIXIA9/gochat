package manager

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/client"
	"sync"
)

const (
	clientsCache     = 200
	roomsCache       = 200
	clientRoomsCache = 200
)

type Manager struct {
	clients     map[domain.UserNumber]client.Client
	rooms       map[domain.RoomNumber]client.Room
	clientRooms map[domain.UserNumber]map[domain.RoomNumber]struct{}
	mu          sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		clients:     make(map[domain.UserNumber]client.Client, clientsCache),
		rooms:       make(map[domain.RoomNumber]client.Room, roomsCache),
		clientRooms: make(map[domain.UserNumber]map[domain.RoomNumber]struct{}, clientRoomsCache),
	}
}

func (m *Manager) Send(ctx context.Context, msg *domain.Message) (bool, error) {
	//todo ctx
	//user
	c, ok := m.clients[domain.UserNumber(msg.To())]
	if ok {
		data, err := msg.MarshalJSON()
		if err != nil {
			return false, err
		}
		c.Write(data)
		return true, nil
	}

	//room
	r, ok := m.rooms[domain.RoomNumber(msg.To())]
	if ok {
		data, err := msg.MarshalJSON()
		if err != nil {
			return false, err
		}

		r.Broadcast(data)
		return true, nil
	}

	return false, nil
}

func (m *Manager) AddClient(c client.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.Number()] = c
}

func (m *Manager) DropClient(number domain.UserNumber) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 关闭并移除 client
	c, exist := m.clients[number]
	if exist {
		if err := c.Close(); err != nil {
			return err
		}
		delete(m.clients, number)
	}

	// 获取并移除 client 的房间集合
	roomNumbers := m.clientRooms[number]
	delete(m.clientRooms, number)

	// 未加入房间直接返回
	if roomNumbers == nil {
		return nil
	}

	// 逐个房间移除 client，并在房间人数 <= 1 时删除房间
	for roomNumber := range roomNumbers {
		if err := m.leaveRoom(number, roomNumber); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) JoinRoom(number domain.RoomNumber, c client.Client) {
	clientNumber := c.Number()
	m.mu.Lock()
	defer m.mu.Unlock()

	room, exist := m.rooms[number]
	if !exist {
		room = client.NewRoom(number)
		m.rooms[number] = room
	}

	room.AddClient(c)

	if m.clientRooms[clientNumber] == nil {
		m.clientRooms[clientNumber] = make(map[domain.RoomNumber]struct{})
	}
	m.clientRooms[clientNumber][number] = struct{}{}
}

func (m *Manager) LeaveRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.leaveRoom(userNumber, roomNumber)
}

func (m *Manager) leaveRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查是否存在
	room, ok := m.rooms[roomNumber]
	if !ok {
		return nil
	}

	//删除房间中的客户端
	room.DropClient(userNumber)
	if rooms, ok := m.clientRooms[userNumber]; ok {
		delete(rooms, roomNumber)

		//删除用户与房间的对应关系
		if len(rooms) == 0 {
			delete(m.clientRooms, userNumber)
		}
	}

	if room.Length() <= 1 {
		delete(m.rooms, roomNumber)
	}
	return nil
}
