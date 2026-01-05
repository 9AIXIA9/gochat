package handler

import (
	"gochat/internal/application"
	websocketInfra "gochat/internal/infrastructure/websocket"
	"gochat/pkg/ctxutil"
	"net/http"

	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type WebsocketHandler struct {
	upgrader                  *websocket.Upgrader
	manager                   *websocketInfra.Manager
	router                    *websocketInfra.Router
	userSessionStartedUseCase application.UserSessionStartedUseCase
}

func NewWebsocketHandler(
	upgrader *websocket.Upgrader,
	manager *websocketInfra.Manager,
	router *websocketInfra.Router,
	userSessionStartedUseCase application.UserSessionStartedUseCase,
) *WebsocketHandler {
	return &WebsocketHandler{
		upgrader:                  upgrader,
		manager:                   manager,
		router:                    router,
		userSessionStartedUseCase: userSessionStartedUseCase,
	}
}

func (s *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := ctxutil.UserIDFrom(r.Context())

	tracer := otel.Tracer("websocket.connection")
	// Short-lived span for the upgrade/handshake so it exports immediately.
	ctxConn, span := tracer.Start(r.Context(), "websocket.upgrade", trace.WithSpanKind(trace.SpanKindServer))
	span.SetAttributes(
		attribute.String("enduser.id", userID.String()),
		attribute.String("net.peer.ip", r.RemoteAddr),
		attribute.String("http.target", r.URL.Path),
	)

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "upgrade failed")
		span.End()
		zap.L().Error(
			"failed to upgrade to websocket",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return
	}
	span.End()

	client := websocketInfra.NewClient(ctxConn, conn, s.router)
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
