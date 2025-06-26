package ctrl

import (
	"gochat/internal/domain"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Manager struct {
	websocket.Upgrader
	roomsMutex sync.RWMutex
	rooms      map[domain.RoomNumber]*Room
}

var (
	manager Manager
)

// JoinRoom 建立连接
func JoinRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber, w http.ResponseWriter, r *http.Request) error {
	//检查是否存在此房间
	room, exist := manager.rooms[roomNumber]
	if !exist {
		return errors.New("the room doesn't exist")
	}

	//检查用户是否已进入
	_, exist = room.clients[userNumber]
	if exist {
		return errors.New("user has existed")
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

// ExitRoom 退出房间
func ExitRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//检查房间是否存在
	room, exist := manager.rooms[roomNumber]
	if !exist {
		return nil
	}

	//关闭用户连接
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
