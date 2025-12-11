package domain

import (
	"time"
)

type RoomProfile struct {
	id           ProfileID
	name         string
	introduction string
	createdAt    time.Time
}

func LoadRoomProfile(
	id ProfileID,
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
	name string,
	introduction string,
	generator ProfileIDGenerator,
) *RoomProfile {
	return &RoomProfile{
		id:           generator.Generate(),
		name:         name,
		introduction: introduction,
		createdAt:    time.Now().UTC(),
	}
}

func (p *RoomProfile) UpdateName(name string) {
	p.name = name
}

func (p *RoomProfile) UpdateIntroduction(introduction string) {
	p.introduction = introduction
}

func (p *RoomProfile) ID() ProfileID {
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
