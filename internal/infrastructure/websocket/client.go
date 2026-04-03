package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/internal/infrastructure/metrics"
	"sync"
	"time"

	myErrors "gochat/internal/shared/errors"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 建议：ping 间隔必须小于读超时（pongWait）
	writeWait  = 2 * time.Second  // 每次写操作超时
	pongWait   = 15 * time.Second // 期待下一次 pong 的最大间隔（读超时）
	pingPeriod = 10 * time.Second // 发送 ping 的间隔，应严格小于 pongWait
)

type Client struct {
	conn   *websocket.Conn
	router *Router

	ctx    context.Context
	cancel context.CancelFunc

	sendChan  chan []byte
	onClose   func()
	closeOnce sync.Once
}

func NewClient(ctx context.Context, conn *websocket.Conn, router *Router) *Client {
	clientCtx, cancel := context.WithCancel(ctx)
	return &Client{
		conn:     conn,
		router:   router,
		ctx:      clientCtx,
		cancel:   cancel,
		sendChan: make(chan []byte, 64),
	}
}

// WithOnClose sets a callback to be invoked exactly once when the client is closed.
func (c *Client) WithOnClose(fn func()) *Client {
	c.onClose = fn
	return c
}

func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
	// 移除额外的 heartbeat 写协程，避免并发写导致连接中断
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.cancel()
		if err := c.conn.Close(); err != nil {
			zap.L().Debug("websocket client close connection failed", zap.Error(err))
		}
		if c.onClose != nil {
			c.onClose()
		}
	})
}

func (c *Client) Send(b []byte) error {
	select {
	case c.sendChan <- b:
		return nil
	default:
		metrics.WSSendFailure(c.ctx, "send_channel_full")
		return myErrors.ErrChanIsFull
	}
}

func (c *Client) readPump() {
	defer c.Close()

	_ = c.conn.SetReadDeadline(time.Now().UTC().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().UTC().Add(pongWait))
	})

	c.conn.SetCloseHandler(func(code int, text string) error {
		zap.L().Debug(
			"websocket client received close frame",
			zap.Int("code", code),
			zap.String("text", text),
		)
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				metrics.WSReadError(c.ctx, "close_frame")
				return
			}
			metrics.WSReadError(c.ctx, "read_failed")
			zap.L().Error("websocket client read message failed", zap.Error(err))
			return
		}

		msg := c.router.Route(c.ctx, message)
		if msg != nil {
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				metrics.WSMessageOut(c.ctx, "route_reply", "marshal_failed")
				zap.L().Error("websocket client marshal message failed", zap.Error(err))
				continue
			}
			if err := c.Send(msgBytes); err != nil {
				metrics.WSMessageOut(c.ctx, "route_reply", "send_failed")
				zap.L().Error("websocket client send message failed", zap.Error(err))
			} else {
				metrics.WSMessageOut(c.ctx, "route_reply", "ok")
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case b := <-c.sendChan:
			_ = c.conn.SetWriteDeadline(time.Now().UTC().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				metrics.WSWriteError(c.ctx, "write_message_failed")
				return
			}
		case <-ticker.C:
			// 统一在单写协程里发送 ping，避免并发写
			_ = c.conn.SetWriteDeadline(time.Now().UTC().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				metrics.WSWriteError(c.ctx, "ping_failed")
				return
			}
		}
	}
}
