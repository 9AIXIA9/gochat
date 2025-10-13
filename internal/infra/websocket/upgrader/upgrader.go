package upgrader

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"time"
)

const bufferSize = 1024

type Upgrader struct {
	upgrader *websocket.Upgrader
}

func NewUpgrader(origins []string) *Upgrader {
	allowed := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		if o != "" {
			allowed[o] = struct{}{}
		}
	}

	u := &websocket.Upgrader{
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

	return &Upgrader{
		upgrader: u,
	}
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return u.upgrader.Upgrade(w, r, responseHeader)
}
