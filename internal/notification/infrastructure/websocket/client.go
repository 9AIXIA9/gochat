package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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
	Wait()
	Close()
}

type client struct {
	conn      *websocket.Conn
	sendChan  chan *Message
	closeOnce sync.Once
	closeCtx  context.Context
	cancelCtx func()
}

func newClient(conn *websocket.Conn) *client {
	ctx, cancel := context.WithCancel(context.Background())

	return &client{
		conn:      conn,
		sendChan:  make(chan *Message, sendCache),
		closeCtx:  ctx,
		cancelCtx: cancel,
	}
}

func (c *client) Start() {
	c.conn.SetReadLimit(readLimit)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		zap.L().Error("set read deadline error", zap.Error(err))
	}

	go c.writePump()
	go c.readPump()
}

// writePump: 处理发送队列并定期发送 ping
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
				if err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed by server")); err != nil {
					zap.L().Error("websocket write close message failed", zap.Error(err))
				}
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
			if err := c.conn.WriteJSON(&Message{
				Type: PingType,
				Data: nil,
			}); err != nil {
				zap.L().Error("websocket json ping failed", zap.Error(err))
				return
			}
		case <-c.closeCtx.Done():
			return
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

		c.handleMsgRead(msg)
	}
}

func (c *client) handleMsgRead(msg *Message) {
	switch msg.Type {
	case PingType:
		// 客户端使用自定义 ping，回复自定义 pong，并延长读超时
		select {
		case c.sendChan <- &Message{Type: PongType}:
		default:
			// 队列满，避免死锁；根据需求选择丢弃
			zap.L().Warn("send channel full on internal message")
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	case PongType:
		// 收到客户端的自定义 pong，延长读超时
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			zap.L().Error("set read deadline failed", zap.Error(err))
		}
	case DataType:
		zap.L().Debug("receive client data message", zap.String("content", string(msg.Data)))
	default:
		zap.L().Error("receive client invalid message type", zap.String("type", string(msg.Type)))
	}
}

func (c *client) Wait() {
	<-c.closeCtx.Done()
}

func (c *client) Close() {
	c.closeOnce.Do(func() {
		c.cancelCtx()
		close(c.sendChan)
		if err := c.conn.Close(); err != nil {
			zap.L().Error("close client failed", zap.Error(err))
		}
	})
}

func (c *client) Send(msg *domain.Message) error {
	data, err := msg.Marshal()
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
		return myErrors.ErrChanIsFull
	}
}
