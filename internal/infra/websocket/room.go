package websocket

import (
	"gochat/internal/domain"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	// 消息广播
	broadcast    chan []byte
	clients      map[domain.UserNumber]Client
	clientsMutex sync.RWMutex
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// func (r *Room) broadcast() {
//
// }
func (r *Room) close() {
	close(r.broadcast)
}

func (c *Client) close() {
	close(c.send)
	if err := c.conn.Close(); err != nil {
		//todo handle error
	}
}

// 接收客户端的信息
func (c *Client) write() {

}

// 发送信息给客户端
func (c *Client) read() {
	//读取历史记录
}
