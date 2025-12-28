package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	maxIntroductionLen = 200
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

func (p *RoomProfile) UpdateName(name string) error {
	if len(name) > maxNameLen {
		return myErrors.NewBusiness("name length exceeds limit: max=%d", maxNameLen)
	}
	p.name = name
	return nil
}

func (p *RoomProfile) UpdateIntroduction(introduction string) error {
	if len(introduction) > maxIntroductionLen {
		return myErrors.NewBusiness("introduction length exceeds limit: max=%d", maxIntroductionLen)
	}
	p.introduction = introduction
	return nil
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
