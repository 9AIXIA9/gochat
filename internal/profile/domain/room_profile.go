package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type RoomProfile struct {
	id           kernel.RoomID
	name         string
	introduction string
	createdAt    time.Time
}

func LoadRoomProfile(
	id kernel.RoomID,
	name string,
	introduction string,
	createdAt time.Time,
) *RoomProfile {
	return &RoomProfile{
		id:           id,
		name:         name,
		introduction: introduction,
		createdAt:    createdAt,
	}
}

func CreateRoomProfile(
	id kernel.RoomID,
	createdAt time.Time,
) *RoomProfile {
	return &RoomProfile{
		id:        id,
		createdAt: createdAt,
	}
}

func (p *RoomProfile) UpdateName(name string) {
	p.name = name
}

func (p *RoomProfile) UpdateIntroduction(introduction string) {
	p.introduction = introduction
}

func (p *RoomProfile) ID() kernel.RoomID {
	return p.id
}

func (p *RoomProfile) Name() string {
	return p.name
}

func (p *RoomProfile) Introduction() string {
	return p.introduction
}

func (p *RoomProfile) CreatedAt() time.Time {
	return p.createdAt
}
