//go:generate mockgen -source=friend_request_repository.go -destination=./mocks/friend_request_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type FriendRequestRepository interface {
	FriendRequestCreator
	FriendRequestUpdater
	FriendRequestFinderByID
	FriendRequestExisterByUserIDAndState
	FriendRequestsFinderByUserID
}

type FriendRequestCreator interface {
	Create(ctx context.Context, friendRequest *FriendRequest) error
}

type FriendRequestUpdater interface {
	Update(ctx context.Context, friendRequest *FriendRequest) error
}

type FriendRequestFinderByID interface {
	FindByID(ctx context.Context, requestID kernel.OperationID) (*FriendRequest, error)
}

type FriendRequestsFinderByUserID interface {
	FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID kernel.OperationID) ([]*FriendRequest, error)
}

type FriendRequestExisterByUserIDAndState interface {
	ExistByUserIDAndState(ctx context.Context, userID1, userID2 kernel.UserID, state FriendRequestState) (bool, error)
}
