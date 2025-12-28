package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type FriendshipID string

type Friendship struct {
	id        FriendshipID
	userID1   kernel.UserID
	userID2   kernel.UserID
	createdAt time.Time

	manager *event.Manager
}

func LoadFriendship(id FriendshipID, userID1, userID2 kernel.UserID, createdAt time.Time) *Friendship {
	return &Friendship{
		id:        id,
		userID1:   userID1,
		userID2:   userID2,
		createdAt: createdAt,
		manager:   event.NewEventManager(),
	}
}

func CreateFriendship(
	userID1, userID2 kernel.UserID,
	friendshipIDGenerator FriendshipIDGenerator,
	idGenerator event.IDGenerator,
) (*Friendship, error) {
	if userID1 == userID2 {
		return nil, ErrAddYourselfAsFriend
	}

	friendship := &Friendship{
		id:        friendshipIDGenerator.Generate(),
		userID1:   userID1,
		userID2:   userID2,
		createdAt: time.Now().UTC(),
		manager:   event.NewEventManager(),
	}

	ev, err := NewFriendshipCreatedEvent(friendship.id, idGenerator)
	if err != nil {
		return nil, err
	}

	friendship.manager.RecordEvent(ev)

	return friendship, nil
}

func (f *Friendship) ID() FriendshipID {
	return f.id
}

func (f *Friendship) UserID1() kernel.UserID {
	return f.userID1
}

func (f *Friendship) UserID2() kernel.UserID {
	return f.userID2
}

func (f *Friendship) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Friendship) GetEvents() []event.Event {
	return f.manager.GetEvents()
}
