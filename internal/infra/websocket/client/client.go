package client

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gochat/internal/domain"
)

const (
	sendCache = 256
	readCache = 256
)

type Client interface {
	Start()
	Write(data []byte)
	Number() domain.UserNumber
	Close() error
}

type client struct {
	number   domain.UserNumber
	conn     *websocket.Conn
	sendChan chan []byte
	readChan chan []byte
}

func NewClient(conn *websocket.Conn, number domain.UserNumber) Client {
	return &client{
		number:   number,
		conn:     conn,
		sendChan: make(chan []byte, sendCache),
		readChan: make(chan []byte, readCache),
	}
}
func (u *client) Start() {
	go u.writePump()
	u.readPump()
}

func (u *client) Number() domain.UserNumber {
	return u.number
}

func (u *client) Write(data []byte) {
	u.sendChan <- data
}

func (u *client) Close() error {
	close(u.sendChan)
	close(u.readChan)
	if err := u.conn.Close(); err != nil {
		return err
	}
	return nil
}

// 发送信息给客户端
func (u *client) writePump() {
	for {
		msg, ok := <-u.sendChan
		if !ok {
			//关闭连接发送的 msg
			if err := u.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				zap.L().Error("websocket write close msg failed", zap.Error(err))
				return
			}
		}

		// 处理发送的 msg
		err := u.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			zap.L().Error("websocket write msg failed", zap.Error(err))
			continue
		}
	}
}

// 接收客户端的信息
func (u *client) readPump() {
	for {
		// 接收客户端的信息
		_, msg, err := u.conn.ReadMessage()
		if err != nil {
			zap.L().Error("websocket read msg failed", zap.Error(err))
			continue
		}

		//todo 客户端发送信息暂时未处理 后续可扩展为已读未读 发送状态等
		u.readChan <- msg
	}
}
