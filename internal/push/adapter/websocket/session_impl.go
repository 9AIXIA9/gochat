package websocket

import (
	"context"
	"gochat/internal/push/contract"
	"gochat/internal/push/core"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait  = 5 * time.Second
	pongWait   = 2 * time.Minute
	pingPeriod = 30 * time.Second
)

// WSSession 是底层真正的 Gorilla WebSocket 实现
type WSSession struct {
	id              string
	userID          string
	conn            *websocket.Conn
	hub             *core.Manager
	upstreamHandler contract.UpstreamHandler
	send            chan []byte // 出发消息缓冲控制 (避免阻塞核心线程)
}

func NewWSSession(id string, userID string, conn *websocket.Conn, hub *core.Manager, handler contract.UpstreamHandler) *WSSession {
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

func (s *WSSession) ID() string {
	return s.id
}

func (s *WSSession) UserID() string {
	return s.userID
}

// Send 这个方法供 Manager 或本地 LocalPushService 调用，纯粹地将数据注入 Channel
func (s *WSSession) Send(msg []byte) error {
	select {
	case s.send <- msg:
		return nil
	default:
		// 当队列满时代表客户端僵死或者处理太慢，按最佳实践抛弃或者断开。
		s.Close()
		return nil // 或者报一个 ErrBufferFull
	}
}

func (s *WSSession) Close() error {
	zap.L().Debug("push adapter: closing websocket session", zap.String("userID", s.userID))
	s.hub.Unregister(s)
	// 清理通道等操作在 writePump defer 内完成会更优雅
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
			s.upstreamHandler.HandleUpstream(ctx, s.userID, message)
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

			// 读尽当前管道内现存积压包做 Batching 支持。
			n := len(s.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'}) // 可以做个定界符，或者由业务包装时保证完整 json
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
