package websocket

import (
	"gochat/internal/gateway/core"
	"gochat/internal/shared/contract"
	"gochat/pkg/ctxutil"
	"net/http"

	"github.com/gorilla/websocket"
)

// IngressHandler 负责承接外网真实流量，建立 WebSocket 握手并接入 gateway 上下文
type IngressHandler struct {
	hub                *core.Manager
	upgrader           websocket.Upgrader
	upstreamHandler    contract.UpstreamHandler
	sessionIDGenerator core.SessionIDGenerator
}

func NewIngressHandler(
	hub *core.Manager,
	handler contract.UpstreamHandler,
	checkOrigin func(r *http.Request) bool,
	idGenerator core.SessionIDGenerator,
) *IngressHandler {
	return &IngressHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     checkOrigin,
		},
		upstreamHandler:    handler,
		sessionIDGenerator: idGenerator,
	}
}

// ServeHTTP 是标准的 http.HandlerFunc 泛型形态方法。
func (h *IngressHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // 降级或升级错误通常交由 Http 框架自身或 upgrader 自带响应，无需额外处理
	}

	session := NewWSSession(h.sessionIDGenerator.Generate(), ctxutil.UserIDFrom(r.Context()), conn, h.hub, h.upstreamHandler)

	// 注册新用户长连接进本地资源池调度
	h.hub.Register(session)

	// 后台开启读写引擎 (pump)
	session.Start(r.Context())
}
