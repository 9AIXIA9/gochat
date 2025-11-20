package gin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
		log.Printf("starting http server on %s", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start http server failed: %s\n", err)
		}
	}()
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server failed: %w", err)
	}
	return nil
}
