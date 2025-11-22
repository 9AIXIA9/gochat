package websocket

import (
	"context"
	"gochat/internal/shared/event"
	"net/http"

	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

//TODO websocket 连接 保持有问题

type Server struct {
	manager *Manager
	router  *Router

	creator     event.UnpublishedEventCreator
	idGenerator event.IDGenerator
}

func NewServer(
	manager *Manager,
	router *Router,
	creator event.UnpublishedEventCreator,
	generator event.IDGenerator,
) *Server {
	return &Server{
		manager:     manager,
		router:      router,
		creator:     creator,
		idGenerator: generator,
	}
}

func (s *Server) ServeWS(w http.ResponseWriter, r *http.Request, userID kernel.UserID) error {
	conn, err := s.manager.Upgrade(w, r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	ctxWithUserID := context.WithValue(ctx, "user_id", userID)
	client := NewClient(ctxWithUserID, conn, s.router)
	// Attach metrics if present in manager
	client.SetMetrics(s.manager.metrics)
	s.manager.Register(userID, client)
	client.Start()

	ev, err := NewUserSessionStartedEvent(s.idGenerator.Generate(), userID)
	if err != nil {
		zap.L().Error("failed to create UserSessionStartedEvent", zap.Error(err))
	}

	if err := s.creator.CreateUnpublishedEvent(ctx, ev); err != nil {
		zap.L().Error("failed to publish UserSessionStartedEvent", zap.Error(err))
	}

	// 当连接关闭时自动注销
	<-ctxWithUserID.Done()
	s.manager.Unregister(userID)

	ev2, err := NewUserSessionEndedEvent(s.idGenerator.Generate(), userID)
	if err != nil {
		zap.L().Error("failed to create UserSessionEndedEvent", zap.Error(err))
	}

	if err := s.creator.CreateUnpublishedEvent(context.Background(), ev2); err != nil {
		zap.L().Error("failed to publish UserSessionEndedEvent", zap.Error(err))
	}

	return nil
}
