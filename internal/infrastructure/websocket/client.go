package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/internal/infrastructure/metrics"
	"gochat/internal/shared/kernel"
	"net"
	"strings"
	"sync"
	"time"

	myErrors "gochat/internal/shared/errors"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 高并发长连接场景下需要更宽松的心跳窗口，避免调度抖动触发误判超时。
	writeWait  = 5 * time.Second
	pongWait   = 2 * time.Minute
	pingPeriod = 30 * time.Second
)

type Client struct {
	id     kernel.UserID
	conn   *websocket.Conn
	router *Router

	ctx    context.Context
	cancel context.CancelFunc

	sendChan  chan []byte
	onClose   func()
	closeOnce sync.Once
}

func NewClient(ctx context.Context, conn *websocket.Conn, router *Router, id kernel.UserID) *Client {
	clientCtx, cancel := context.WithCancel(ctx)
	return &Client{
		id:       id,
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

	refreshReadDeadline := func() error {
		return c.conn.SetReadDeadline(time.Now().UTC().Add(pongWait))
	}

	_ = refreshReadDeadline()
	c.conn.SetPongHandler(func(string) error {
		return refreshReadDeadline()
	})
	c.conn.SetPingHandler(func(appData string) error {
		if err := refreshReadDeadline(); err != nil {
			return err
		}
		_ = c.conn.SetWriteDeadline(time.Now().UTC().Add(writeWait))
		return c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().UTC().Add(writeWait))
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
			if c.ctx.Err() != nil || isExpectedReadCloseError(err) {
				metrics.WSReadError(c.ctx, "connection_closed")
				return
			}
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				metrics.WSReadError(c.ctx, "read_timeout")
				zap.L().Debug("websocket client read timeout", zap.Error(err))
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

func isExpectedReadCloseError(err error) bool {
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "forcibly closed by the remote host") ||
		strings.Contains(errMsg, "use of closed network connection") ||
		strings.Contains(errMsg, "connection reset by peer")
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
