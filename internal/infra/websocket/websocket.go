package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gochat/internal/domain"
	"gochat/internal/types"
	"net/http"
	"sync"
)

//todo 事件驱动模式

type Manager struct {
	websocket.Upgrader
	roomsMutex sync.RWMutex
	rooms      map[domain.RoomNumber]*Room
}

var (
	//todo 不要包级全局变量
	manager Manager
)

// CreateGroup 创建房间
func CreateGroup(owner domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查是否存在此房间
	manager.roomsMutex.Lock()
	defer manager.roomsMutex.Unlock()

	_, exist := manager.rooms[roomNumber]
	if exist {
		return types.ErrRoomExist
	}

	//创建房间
	room := &Room{
		broadcast:    make(chan []byte, 256),
		clients:      nil,
		clientsMutex: sync.RWMutex{},
	}

	manager.rooms[roomNumber] = room
	return nil
}

// EstablishConnection 建立连接
func EstablishConnection(w http.ResponseWriter, r *http.Request, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查是否存在此房间
	manager.roomsMutex.RLock()
	defer manager.roomsMutex.RUnlock()

	room, exist := manager.rooms[roomNumber]
	if !exist {
		return types.ErrRoomNotExist
	}

	//检查用户是否已进入
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	_, exist = room.clients[userNumber]
	if exist {
		return types.ErrUserExist
	}

	// 升级HTTP连接到WebSocket
	ws, err := manager.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("upgrade error:%w", err)
	}

	//创建客户端
	client := Client{
		conn: ws,
		send: make(chan []byte, 256),
	}

	room.clients[userNumber] = client

	//持续接收消息
	go client.read()
	go client.write()
	return nil
}

// ExitConnection 关闭连接
func ExitConnection(userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查房间是否存在
	manager.roomsMutex.Lock()
	defer manager.roomsMutex.Unlock()

	room, exist := manager.rooms[roomNumber]
	if !exist {
		return nil
	}

	//关闭用户连接
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	client, exist := room.clients[userNumber]
	if exist {
		client.close()
		delete(room.clients, userNumber)

		if len(room.clients) == 0 {
			room.close()
			delete(manager.rooms, roomNumber)
		}
	}
	return nil
}
