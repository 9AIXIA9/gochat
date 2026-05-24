package gin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	Router      *gin.Engine
	Server      *http.Server
	DebugServer *http.Server
}

func NewServer(router *gin.Engine, server *http.Server, debugServer *http.Server) *Server {
	return &Server{
		Router:      router,
		Server:      server,
		DebugServer: debugServer,
	}
}

func (s *Server) Start() {
	startServer := func(server *http.Server, label string) {
		if server == nil {
			return
		}
		go func() {
			zap.L().Info("starting http server", zap.String("label", label), zap.String("addr", server.Addr))
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Error("start http server failed", zap.String("label", label), zap.Error(err))
			}
		}()
	}

	startServer(s.Server, "api")
	startServer(s.DebugServer, "pprof")
}

func (s *Server) Close(ctx context.Context) error {
	if s.DebugServer != nil {
		if err := s.DebugServer.Shutdown(ctx); err != nil {
			zap.L().Error("shutdown pprof server failed", zap.Error(err))
		}
	}

	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server failed: %w", err)
	}
	return nil
}

func NewPProfMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return mux
}
