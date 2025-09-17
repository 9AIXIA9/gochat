package client

import "gochat/internal/domain"

const (
	roomCache = 10
)

type Room interface {
	AddClient(c Client)
	DropClient(number domain.UserNumber)
	Number() domain.RoomNumber
	Broadcast(data []byte)
	Length() int
}

type room struct {
	number  domain.RoomNumber
	clients map[domain.UserNumber]Client
}

func NewRoom(number domain.RoomNumber) Room {
	return &room{
		number:  number,
		clients: make(map[domain.UserNumber]Client, roomCache),
	}
}

func (r *room) AddClient(c Client) {
	r.clients[c.Number()] = c
}

func (r *room) DropClient(number domain.UserNumber) {
	delete(r.clients, number)
}

func (r *room) Number() domain.RoomNumber {
	return r.number
}

func (r *room) Broadcast(data []byte) {
	for _, c := range r.clients {
		c.Write(data)
	}
}

func (r *room) Length() int {
	return r.Length()
}
