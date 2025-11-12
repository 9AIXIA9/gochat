package websocket

import (
	"net/http"

	"gochat/internal/shared/kernel"
)

type Server struct {
	Manager *Manager
	Router  *Router
}

func NewServer(m *Manager, r *Router) *Server {
	return &Server{
		Manager: m,
		Router:  r,
	}
}

func (s *Server) ServeWS(w http.ResponseWriter, r *http.Request, userID kernel.UserID) error {
	conn, err := s.Manager.Upgrade(w, r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	client := NewClient(ctx, conn, s.Router)
	s.Manager.Register(userID, client)
	client.Start()

	// 当连接关闭时自动注销
	<-ctx.Done()
	s.Manager.Unregister(userID)
	return nil
}
