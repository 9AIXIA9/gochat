package gin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	Router *gin.Engine
	Server *http.Server
}

func NewServer(router *gin.Engine, server *http.Server) *Server {
	return &Server{
		Router: router,
		Server: server,
	}
}

func (s *Server) Start() {
	// 启动 HTTP 服务器
	go func() {
		zap.L().Info("starting http server", zap.String("addr", s.Server.Addr))
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("start http server failed", zap.Error(err))
		}
	}()
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server failed: %w", err)
	}
	return nil
}
