package upgrader

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

const bufferSize = 1024

var upgrader websocket.Upgrader
var once sync.Once

func Init(origins []string) {
	once.Do(func() {
		allowed := make(map[string]struct{}, len(origins))
		for _, o := range origins {
			if o != "" {
				allowed[o] = struct{}{}
			}
		}

		upgrader = websocket.Upgrader{
			HandshakeTimeout: 5 * time.Second,
			ReadBufferSize:   bufferSize,
			WriteBufferSize:  bufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
				//// 若未配置 origins 则允许所有来源
				//if len(allowed) == 0 {
				//	return true
				//}
				//origin := r.Header.Get("Origin")
				//if origin == "" {
				//	return false
				//}
				//
				//_, ok := allowed[origin]
				//return ok
			},
			EnableCompression: false,
		}
	})
}

func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, responseHeader)
}
