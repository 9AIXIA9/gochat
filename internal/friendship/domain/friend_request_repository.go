//go:generate mockgen -source=friend_request_repository.go -destination=./mocks/friend_request_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type FriendRequestRepository interface {
	FriendRequestExister
}

type FriendRequestExister interface {
	Exist(ctx context.Context, from kernel.UserID, to kernel.UserID) (bool, error)
}
