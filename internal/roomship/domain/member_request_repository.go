//go:generate mockgen -source=member_request_repository.go -destination=./mocks/mock_member_request_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type MemberRequestRepository interface {
	MemberRequestCreator
	MemberRequestUpdater
	MemberRequestFinderByID
	MemberRequestExisterByUserIDAndRoomIDAndState
	MemberRequestsFinderByUserID
}

type MemberRequestCreator interface {
	Create(ctx context.Context, memberRequest *MemberRequest) error
}

type MemberRequestUpdater interface {
	Update(ctx context.Context, memberRequest *MemberRequest) error
}

type MemberRequestFinderByID interface {
	FindByID(ctx context.Context, requestID kernel.OperationID) (*MemberRequest, error)
}

type MemberRequestExisterByUserIDAndRoomIDAndState interface {
	ExistByUserIDAndRoomIDAndState(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID, state MemberRequestState) (bool, error)
}

type MemberRequestsFinderByUserID interface {
	FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID kernel.OperationID) ([]*MemberRequest, error)
}
