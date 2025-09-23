package client

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/types"
	"sync"
	"sync/atomic"
)

const (
	sendCache = 256
)

type Client interface {
	Start()
	Write(data []byte) error
	Number() domain.UserNumber
	Close() error
}

// todo 心跳检测
type client struct {
	number    domain.UserNumber
	conn      *websocket.Conn
	sendChan  chan []byte
	onClose   func(number domain.UserNumber)
	closeOnce sync.Once
	closed    atomic.Bool
	closeMu   sync.RWMutex
}

func NewClient(conn *websocket.Conn, number domain.UserNumber, closeFn func(number domain.UserNumber)) Client {
	return &client{
		number:   number,
		conn:     conn,
		sendChan: make(chan []byte, sendCache),
		onClose:  closeFn,
	}
}

func (c *client) Start() {
	go c.writePump()
	c.readPump()
}

func (c *client) Number() domain.UserNumber {
	return c.number
}

func (c *client) Write(data []byte) error {
	c.closeMu.RLock()
	defer c.closeMu.RUnlock()

	if !c.closed.Load() {
		select {
		case c.sendChan <- data:
			return nil
		default:
			return types.ErrFullMessage
		}
	}
	return types.ErrClientClosed
}

func (c *client) Close() (err error) {
	c.closeOnce.Do(func() {
		c.closeMu.Lock()
		defer c.closeMu.Unlock()

		c.closed.CompareAndSwap(false, true)
		close(c.sendChan)

		//关闭连接
		if err2 := c.conn.Close(); err2 != nil {
			err = err2
		}

		//调用关闭回调
		if c.onClose != nil {
			c.onClose(c.number)
		}
	})
	return err
}

// 发送信息给客户端
func (c *client) writePump() {
	defer func() {
		//关闭客户端
		if err := c.Close(); err != nil {
			zap.L().Error("close client failed", zap.Error(err))
		}
	}()

	for bytes := range c.sendChan {
		// 处理发送的 msg
		err := c.conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			zap.L().Error("websocket write msg failed", zap.Error(err))
			return
		}
	}
}

// 接收客户端的信息
func (c *client) readPump() {
	for {
		// 接收客户端的信息
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			zap.L().Info("websocket client closed", zap.Error(err))
			err := c.Close()
			if err != nil {
				zap.L().Error("websocket client close failed", zap.Error(err))
				return
			}
			return
		}

		//todo 客户端发送信息暂时未处理 后续可扩展为已读未读 ASK确认
		zap.L().Info("receive message sent from client",
			zap.String("content", string(msg)),
			zap.Int64("number", int64(c.Number())))
	}
}
