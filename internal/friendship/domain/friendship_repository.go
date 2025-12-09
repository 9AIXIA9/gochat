//go:generate mockgen -source=friendship_repository.go -destination=./mocks/friendship_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type FriendshipRepository interface {
	FriendshipCreator
	FriendshipFinderByID
	FriendshipExisterByUserID
	FriendshipsFinderByUserID
}

type FriendshipCreator interface {
	Create(ctx context.Context, friendship *Friendship) error
}

type FriendshipFinderByID interface {
	FindByID(ctx context.Context, friendshipID FriendshipID) (*Friendship, error)
}

type FriendshipExisterByUserID interface {
	ExistByUserID(ctx context.Context, userID1, userID2 kernel.UserID) (bool, error)
}

type FriendshipsFinderByUserID interface {
	FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID FriendshipID) ([]*Friendship, error)
}
