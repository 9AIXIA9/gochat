package websocket

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const bufferSize = 1024

func NewUpgrader(origins []string) *websocket.Upgrader {
	allowed := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		if o != "" {
			allowed[o] = struct{}{}
		}
	}

	return &websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   bufferSize,
		WriteBufferSize:  bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			// 若未配置 origins 则允许所有来源
			if len(allowed) == 0 {
				return true
			}

			from := r.RemoteAddr
			if strings.Contains(from, "127.0.0.1") || strings.Contains(from, "localhost") {
				return true
			}

			origin := r.Header.Get("Origin")
			if origin == "" {
				return false
			}

			_, ok := allowed[origin]
			return ok
		},
		EnableCompression: false,
	}
}
