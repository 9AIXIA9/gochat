package websocket

import (
	"context"
	"errors"
	"time"

	"gochat/internal/infrastructure/prometheus"
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

	send    chan []byte
	metrics *prometheus.Metrics
}

func NewClient(ctx context.Context, conn *websocket.Conn, router *Router) *Client {
	clientCtx, cancel := context.WithCancel(ctx)
	return &Client{
		conn:   conn,
		router: router,
		ctx:    clientCtx,
		cancel: cancel,
		send:   make(chan []byte, 64),
	}
}

func (c *Client) SetMetrics(m *prometheus.Metrics) { c.metrics = m }

func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}

func (c *Client) Close() {
	c.cancel()
	if err := c.conn.Close(); err != nil {
		zap.L().Debug("websocket client close connection failed", zap.Error(err))
		return
	}
}

func (c *Client) Send(b []byte) error {
	select {
	case c.send <- b:
		return nil
	default:
		return myErrors.ErrChanIsFull
	}
}

func (c *Client) readPump() {
	defer c.Close()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
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
				if c.metrics != nil {
					c.metrics.WSReadErrors.WithLabelValues("close").Inc()
				}
				return
			}
			if c.metrics != nil {
				c.metrics.WSReadErrors.WithLabelValues("other").Inc()
			}
			return
		}

		msg, err := DecodeMessage(message)
		if err != nil {
			continue
		}

		switch msg.Type {
		case PingType:
			if c.metrics != nil {
				c.metrics.WSMessagesIn.WithLabelValues("ping").Inc()
			}
			if b, err := EncodePong(); err == nil {
				_ = c.Send(b)
				if c.metrics != nil {
					c.metrics.WSMessagesOut.WithLabelValues("pong").Inc()
				}
			}
		case RequestType:
			if c.metrics != nil {
				c.metrics.WSMessagesIn.WithLabelValues("request").Inc()
			}
			req, err := DecodeRequest(msg.Payload)
			if err != nil {
				zap.L().Debug("websocket decode request failed", zap.Error(err))
				continue
			}
			start := time.Now()
			resp, err := c.router.Route(c.ctx, req)
			if c.metrics != nil {
				c.metrics.WSRouteDur.WithLabelValues(string(req.RequestTopic)).Observe(time.Since(start).Seconds())
			}
			if err != nil {
				zap.L().Error("websocket route request failed", zap.Error(err))
				continue
			}

			if resp != nil {
				b, err := EncodeResponse(resp)
				if err != nil {
					zap.L().Debug("websocket encode response failed", zap.Error(err))
					continue
				}
				_ = c.Send(b)
				if c.metrics != nil {
					c.metrics.WSMessagesOut.WithLabelValues("response").Inc()
				}
			}
		default:
			if c.metrics != nil {
				c.metrics.WSMessagesIn.WithLabelValues("other").Inc()
			}
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
		case b := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				if c.metrics != nil {
					c.metrics.WSWriteErrors.WithLabelValues("text").Inc()
				}
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if c.metrics != nil {
					c.metrics.WSWriteErrors.WithLabelValues("ping").Inc()
				}
				return
			}
		}
	}
}
