package domain

import "context"

// UnitOfWork 工作单元
type UnitOfWork[T any] func(T) error

type JoinRoomUnitOfWork UnitOfWork[JoinRoomAggregate]

type JoinRoomAggregateUOW interface {
	DoJoinRoomUOW(ctx context.Context, uow JoinRoomUnitOfWork) error
	JoinRoomAggregate
}
