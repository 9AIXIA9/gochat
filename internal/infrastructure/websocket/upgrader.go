package websocket

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const bufferSize = 1024

func NewUpgrader(origins []string) *websocket.Upgrader {
	allowed := make(map[string]struct{}, len(origins))
	allowAll := false
	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			allowAll = true
			continue
		}
		allowed[o] = struct{}{}
	}

	return &websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   bufferSize,
		WriteBufferSize:  bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			// 若未配置 origins 或包含通配符则允许所有来源
			if allowAll || len(allowed) == 0 {
				return true
			}

			if isLoopbackRequest(r.RemoteAddr) {
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

func isLoopbackRequest(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	host = strings.Trim(host, "[]")
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
