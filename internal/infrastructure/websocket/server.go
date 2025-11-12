package websocket

import (
	"gochat/internal/shared/event"
	"net/http"

	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

type Server struct {
	manager *Manager
	router  *Router

	publisher   event.Publisher
	idGenerator event.IDGenerator
}

func NewServer(
	manager *Manager,
	router *Router,
	publisher event.Publisher,
	generator event.IDGenerator,
) *Server {
	return &Server{
		manager:     manager,
		router:      router,
		publisher:   publisher,
		idGenerator: generator,
	}
}

func (s *Server) ServeWS(w http.ResponseWriter, r *http.Request, userID kernel.UserID) error {
	conn, err := s.manager.Upgrade(w, r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	client := NewClient(ctx, conn, s.router)
	s.manager.Register(userID, client)
	client.Start()

	ev, err := NewUserSessionStartedEvent(s.idGenerator.Generate(), userID)
	if err != nil {
		zap.L().Error("failed to create UserSessionStartedEvent", zap.Error(err))
	}

	if err := s.publisher.Publish(ev); err != nil {
		zap.L().Error("failed to publish UserSessionStartedEvent", zap.Error(err))
	}

	// 当连接关闭时自动注销
	<-ctx.Done()
	s.manager.Unregister(userID)

	ev2, err := NewUserSessionEndedEvent(s.idGenerator.Generate(), userID)
	if err != nil {
		zap.L().Error("failed to create UserSessionEndedEvent", zap.Error(err))
	}

	if err := s.publisher.Publish(ev2); err != nil {
		zap.L().Error("failed to publish UserSessionEndedEvent", zap.Error(err))
	}

	return nil
}
