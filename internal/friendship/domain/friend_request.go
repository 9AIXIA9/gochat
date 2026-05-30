package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type FriendRequestState string

func (s FriendRequestState) String() string {
	return string(s)
}

const (
	StatePending FriendRequestState = "pending"
	StateAgreed  FriendRequestState = "agreed"
	StateRefused FriendRequestState = "refused"
)

const maxContentLength = 100

type FriendRequest struct {
	id      kernel.OperationID
	from    kernel.UserID
	to      kernel.UserID
	content string
	state   FriendRequestState
	sentAt  time.Time //UTC

	eventManager *event.Manager
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
		id:           id,
		from:         from,
		to:           to,
		content:      content,
		state:        state,
		sentAt:       sentAt,
		eventManager: event.NewEventManager(),
	}
}

func CreateFriendRequest(
	from kernel.UserID,
	to kernel.UserID,
	content string,
	operationIDGenerator kernel.OperationIDGenerator,
) (*FriendRequest, error) {
	if from == to {
		return nil, ErrAddYourselfAsFriend
	}

	if len(content) > maxContentLength {
		return nil, ErrFriendRequestContentTooLong
	}

	req := &FriendRequest{
		id:           operationIDGenerator.Generate(),
		from:         from,
		to:           to,
		content:      content,
		state:        StatePending,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}
	return req, nil
}

func (r *FriendRequest) Agree(
	idGenerator event.IDGenerator,
) error {
	if r.state != StatePending {
		return ErrFriendRequestHasBeenHandled
	}

	r.state = StateAgreed

	ev, err := NewFriendRequestAgreedEvent(r.id, idGenerator)
	if err != nil {
		return err
	}

	r.eventManager.RecordEvent(ev)
	return nil
}

func (r *FriendRequest) Refuse() error {
	if r.state != StatePending {
		return ErrFriendRequestHasBeenHandled
	}

	r.state = StateRefused
	return nil
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

func (r *FriendRequest) GetEvents() []event.Event {
	return r.eventManager.GetEvents()
}
