package websocket

import (
	"gochat/internal/application"
	"gochat/pkg/utils"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	upgrader                  *websocket.Upgrader
	manager                   *Manager
	router                    *Router
	userSessionStartedUseCase application.UserSessionStartedUseCase
}

func NewServer(
	upgrader *websocket.Upgrader,
	manager *Manager,
	router *Router,
	userSessionStartedUseCase application.UserSessionStartedUseCase,
) *Server {
	return &Server{
		upgrader:                  upgrader,
		manager:                   manager,
		router:                    router,
		userSessionStartedUseCase: userSessionStartedUseCase,
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
	client.WithOnClose(func() {
		// Unregister by pointer to avoid removing a newly registered client when replacing connections.
		s.manager.UnregisterClient(client)
	})

	s.manager.Register(userID, client)
	client.Start()

	if _, err := s.userSessionStartedUseCase.Execute(r.Context(), &application.UserSessionStartedInput{
		UserID: userID,
	}); err != nil {
		zap.L().Debug(
			"failed to execute user session started use case",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
	}

	<-r.Context().Done()
}
