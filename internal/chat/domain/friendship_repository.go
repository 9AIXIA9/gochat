//go:generate mockgen -source=friendship_repository.go -destination=./mocks/mock_friendship_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type FriendshipRepository interface {
	FriendshipCreator
	FriendshipExisterByUserID
}

type FriendshipCreator interface {
	Create(ctx context.Context, friendship *Friendship) error
}

type FriendshipExisterByUserID interface {
	ExistByUserID(ctx context.Context, userID1, userID2 kernel.UserID) (bool, error)
}
