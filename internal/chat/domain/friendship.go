package domain

import "gochat/internal/shared/kernel"

type FriendshipID string

func (id FriendshipID) String() string {
	return string(id)
}

type Friendship struct {
	id      FriendshipID
	userID1 kernel.UserID
	userID2 kernel.UserID
}

func LoadFriendship(
	id FriendshipID,
	userID1 kernel.UserID,
	userID2 kernel.UserID,
) *Friendship {
	return &Friendship{
		id:      id,
		userID1: userID1,
		userID2: userID2,
	}
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
