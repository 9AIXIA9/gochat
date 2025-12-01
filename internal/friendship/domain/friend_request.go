package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type FriendRequestState string

func (s FriendRequestState) String() string {
	return string(s)
}

const (
	StatePending FriendRequestState = "pending"
	StateAgree   FriendRequestState = "agree"
	StateRefused FriendRequestState = "refused"
)

type FriendRequest struct {
	id      kernel.OperationID
	from    kernel.UserID
	to      kernel.UserID
	content string
	state   FriendRequestState
	sentAt  time.Time //UTC
}

func LoadFriendRequest(
	id kernel.OperationID,
	from kernel.UserID,
	to kernel.UserID,
	content string,
	state FriendRequestState,
	sentAt time.Time,
) *FriendRequest {
	return &FriendRequest{
		id:      id,
		from:    from,
		to:      to,
		content: content,
		state:   state,
		sentAt:  sentAt,
	}
}

func CreateFriendRequest(
	from kernel.UserID,
	to kernel.UserID,
	content string,
	idGenerator OperationIDGenerator,
) *FriendRequest {
	return &FriendRequest{
		id:      idGenerator.Generate(),
		from:    from,
		to:      to,
		content: content,
		state:   StatePending,
		sentAt:  time.Now().UTC(),
	}
}

func (r *FriendRequest) ID() kernel.OperationID {
	return r.id
}

func (r *FriendRequest) From() kernel.UserID {
	return r.from
}

func (r *FriendRequest) To() kernel.UserID {
	return r.to
}

func (r *FriendRequest) Content() string {
	return r.content
}

func (r *FriendRequest) SentAt() time.Time {
	return r.sentAt
}

func (r *FriendRequest) State() FriendRequestState {
	return r.state
}
