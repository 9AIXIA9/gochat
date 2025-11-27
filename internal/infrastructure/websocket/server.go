package websocket

import (
	"gochat/pkg/utils"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	upgrader *websocket.Upgrader
	manager  *Manager
	router   *Router
}

func NewServer(
	upgrader *websocket.Upgrader,
	manager *Manager,
	router *Router,
) *Server {
	return &Server{
		upgrader: upgrader,
		manager:  manager,
		router:   router,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r.Context())

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.L().Error(
			"failed to upgrade to websocket",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return
	}

	client := NewClient(r.Context(), conn, s.router)
	// Attach metrics if present in manager
	s.manager.Register(userID, client)
	client.Start()

	// 当连接关闭时自动注销
	<-r.Context().Done()
	s.manager.Unregister(userID)
}
