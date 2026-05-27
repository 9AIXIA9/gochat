package handler

import (
	"gochat/internal/infrastructure/metrics"
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
	upgrader *websocket.Upgrader
	manager  *websocketInfra.Manager
	router   *websocketInfra.Router
}

func NewWebsocketHandler(
	upgrader *websocket.Upgrader,
	manager *websocketInfra.Manager,
	router *websocketInfra.Router,
) *WebsocketHandler {
	return &WebsocketHandler{
		upgrader: upgrader,
		manager:  manager,
		router:   router,
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
		metrics.WSHandshake(r.Context(), "upgrade_failed")
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
	metrics.WSHandshake(r.Context(), "ok")
	span.End()

	client := websocketInfra.NewClient(ctxConn, conn, s.router, userID)
	client.WithOnClose(func() {
		// Unregister by pointer to avoid removing a newly registered client when replacing connections.
		s.manager.Unregister(client)
	})

	s.manager.Register(userID, client)
	client.Start()

	<-r.Context().Done()
}
