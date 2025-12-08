package websocket

import (
	"context"
	"errors"
	"sync"
	"time"

	myErrors "gochat/internal/shared/errors"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
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

func (c *Client) SendResponse(resp *Response) error {
	b, err := EncodeResponse(resp)
	if err != nil {
		return err
	}
	return c.send(b)
}

func (c *Client) send(b []byte) error {
	select {
	case c.sendChan <- b:
		return nil
	default:
		return myErrors.ErrChanIsFull
	}
}

func (c *Client) readPump() {
	defer c.Close()

	_ = c.conn.SetReadDeadline(time.Now().UTC().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().UTC().Add(pongWait))
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
				return
			}
			zap.L().Error("websocket client read message failed", zap.Error(err))
			return
		}

		msg, err := DecodeMessage(message)
		if err != nil {
			zap.L().Debug("websocket decode message failed", zap.Error(err))
			continue
		}

		switch msg.Type {
		case PingType:
			if b, err := EncodePong(); err == nil {
				_ = c.send(b)
			}
		case RequestType:
			req, err := DecodeRequest(msg.Payload)
			if err != nil {
				zap.L().Debug("websocket decode request failed", zap.Error(err))
				continue
			}
			resp := c.router.Route(c.ctx, req)

			if resp != nil {
				_ = c.SendResponse(resp)
			}
		default:
			zap.L().Debug("websocket client receive invalid message type", zap.String("type", string(msg.Type)))
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
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().UTC().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
