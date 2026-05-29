package websocket

import (
	"context"
	"gochat/internal/gateway/core"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait  = 5 * time.Second
	pongWait   = 2 * time.Minute
	pingPeriod = 30 * time.Second
)

type WSSession struct {
	id              core.SessionID
	userID          kernel.UserID
	conn            *websocket.Conn
	hub             *core.Manager
	upstreamHandler contract.UpstreamHandler
	send            chan []byte
}

func NewWSSession(
	id core.SessionID,
	userID kernel.UserID,
	conn *websocket.Conn,
	hub *core.Manager,
	handler contract.UpstreamHandler,
) *WSSession {
	return &WSSession{
		id:              id,
		userID:          userID,
		conn:            conn,
		hub:             hub,
		upstreamHandler: handler,
		send:            make(chan []byte, 256),
	}
}

// Start 启动当前连接的收发 pump Goroutine
func (s *WSSession) Start(ctx context.Context) {
	go s.readPump(ctx)
	go s.writePump(ctx)
}

func (s *WSSession) ID() core.SessionID {
	return s.id
}

func (s *WSSession) UserID() kernel.UserID {
	return s.userID
}

func (s *WSSession) Send(msg []byte) error {
	select {
	case s.send <- msg:
		return nil
	default:
		s.Close()
		return nil
	}
}

func (s *WSSession) Close() error {
	zap.L().Debug("gateway adapter: closing websocket session", zap.String("userID", s.userID.String()))
	s.hub.Unregister(s)
	return s.conn.Close()
}

func (s *WSSession) readPump(ctx context.Context) {
	defer func() {
		s.Close()
	}()
	s.conn.SetReadLimit(4096)
	_ = s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error {
		_ = s.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.L().Warn("websocket read abnormal error", zap.Error(err))
			}
			break
		}

		// 读到了上行消息 => 立马丢给 UpstreamHandler，自身不含任何业务逻辑，纯透传。
		if s.upstreamHandler != nil {
			if ackMsg := s.upstreamHandler.HandleUpstream(ctx, s.userID, message); ackMsg != nil {
				// 直接塞入当前连接的写队列
				s.Send(ackMsg)
			}
		}
	}
}

func (s *WSSession) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.Close()
	}()

	for {
		select {
		case message, ok := <-s.send:
			_ = s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// 读尽当前管道内现存积压包做 Batching 支持
			n := len(s.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-s.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
