package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/types"
	"sync"
	"time"
)

const (
	sendCache = 256
	// 控制帧与超时相关（仍用这些常量，但不使用控制帧）
	writeWait  = 10 * time.Second
	pongWait   = 90 * time.Second
	pingPeriod = (pongWait * 9) / 10
	readLimit  = 500
)

type Client interface {
	Start()
	Write(msg *domain.Message) error
	Number() domain.UserNumber
	SetCloseHandler(func(code int, text string) error)
	Close()
}

type client struct {
	number    domain.UserNumber
	conn      *websocket.Conn
	sendChan  chan *Message
	closeMu   sync.RWMutex
	closeOnce sync.Once
}

func New(conn *websocket.Conn, number domain.UserNumber) Client {
	return &client{
		number:   number,
		conn:     conn,
		sendChan: make(chan *Message, sendCache),
	}
}

func (c *client) Start() {
	c.conn.SetReadLimit(readLimit)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	go c.writePump()
	c.readPump()
}

func (c *client) Number() domain.UserNumber {
	return c.number
}

func (c *client) Write(msg *domain.Message) error {
	c.closeMu.RLock()
	defer c.closeMu.RUnlock()

	data, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	select {
	case c.sendChan <- &Message{
		Type: DataType,
		Data: data,
	}:
		return nil
	default:
		return types.ErrFullMessage
	}
}

// writePump: 处理发送队列并定期发送自定义 JSON ping
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-c.sendChan:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				zap.L().Error("set write deadline failed", zap.Error(err))
				return
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				zap.L().Error("websocket write msg failed", zap.Error(err))
				return
			}
		case <-ticker.C:
			// 发送自定义 JSON ping（前端会以 JSON pong 回复）
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				zap.L().Error("set write deadline failed", zap.Error(err))
				return
			}
			if err := c.conn.WriteJSON(&Message{Type: PingType}); err != nil {
				zap.L().Error("websocket json ping failed", zap.Error(err))
				return
			}
		}
	}
}

// readPump: 读取客户端消息，收到自定义 pong 时手动延长 ReadDeadline
func (c *client) readPump() {
	for {
		_, rawData, err := c.conn.ReadMessage()
		if err != nil {
			zap.L().Info("websocket client closed", zap.Error(err))
			c.Close()
			return
		}

		msg := new(Message)
		if err := json.Unmarshal(rawData, msg); err != nil {
			zap.L().Error("websocket client receive invalid message", zap.Error(err))
			continue
		}

		c.handleClientMsg(msg)
	}
}

func (c *client) handleClientMsg(msg *Message) {
	switch msg.Type {
	case PingType:
		// 客户端使用自定义 ping，回复自定义 pong，并延长读超时
		if err := c.conn.WriteJSON(&Message{Type: PongType}); err != nil {
			zap.L().Error("response client pong failed", zap.Error(err))
			c.Close()
			return
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	case PongType:
		// 收到客户端的自定义 pong，延长读超时
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	case DataType:
		zap.L().Debug("receive client data message", zap.String("content", string(msg.Data)))
	default:
		zap.L().Error("receive client invalid message type", zap.String("type", string(msg.Type)))
	}
}

func (c *client) Close() {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()

	c.closeOnce.Do(func() {
		close(c.sendChan)
		if err := c.conn.Close(); err != nil {
			zap.L().Error("close client failed", zap.Error(err))
		}
	})
}

func (c *client) SetCloseHandler(fn func(code int, text string) error) {
	c.conn.SetCloseHandler(fn)
}
