package websocket

import (
	"go.uber.org/zap"
	"gochat/internal/domain"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	// 消息广播
	broadcastChan chan []byte
	clients       map[domain.UserNumber]Client
	clientsMutex  sync.RWMutex
}

type Client struct {
	conn      *websocket.Conn
	sendChan  chan []byte
	broadcast chan []byte
}

func (r *Room) broadcast() {
	for {
		msg := <-r.broadcastChan
		r.clientsMutex.RLock()
		for _, client := range r.clients {
			client.sendChan <- msg
		}
		r.clientsMutex.RUnlock()
	}
}

func (r *Room) close() {
	close(r.broadcastChan)
}

// 发送信息给客户端
func (c *Client) writePump() {
	for {
		msg, ok := <-c.sendChan
		if !ok {
			//关闭 msg
			if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				zap.L().Error("client connection write close msg failed", zap.Error(err))
				return
			}
		}

		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			zap.L().Error("client connection write msg failed", zap.Error(err))
			continue
		}
	}
}

// 接收客户端的信息
func (c *Client) readPump() {
	//todo 读取历史记录
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			zap.L().Error("client connection read msg failed", zap.Error(err))
			continue
		}
		c.broadcast <- msg
	}
}

func (c *Client) close() error {
	close(c.sendChan)
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}
