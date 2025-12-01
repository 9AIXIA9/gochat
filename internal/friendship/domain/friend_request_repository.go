//go:generate mockgen -source=friend_request_repository.go -destination=./mocks/friend_request_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type FriendRequestRepository interface {
	FriendRequestsFinderByUserID
}

type FriendRequestsFinderByUserID interface {
	FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID ...kernel.OperationID) ([]*FriendRequest, error)
}
