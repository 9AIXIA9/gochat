package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type User struct {
	id     kernel.UserID
	number kernel.UserNumber

	friendIDs    []kernel.UserID
	requests     []*FriendRequest
	eventManager *event.Manager
}

func LoadUser(
	id kernel.UserID,
	number kernel.UserNumber,
	friendIDs []kernel.UserID,
	requests []*FriendRequest,
) *User {
	return &User{
		id:           id,
		number:       number,
		friendIDs:    friendIDs,
		requests:     requests,
		eventManager: event.NewEventManager(),
	}
}

func CreateUser(
	id kernel.UserID,
	number kernel.UserNumber,
) *User {
	return &User{
		id:           id,
		number:       number,
		friendIDs:    make([]kernel.UserID, 0),
		requests:     make([]*FriendRequest, 0),
		eventManager: event.NewEventManager(),
	}
}

func (u *User) SendFriendRequest(
	to kernel.UserID,
	content string,
	operationIDGenerator OperationIDGenerator,
) (*FriendRequest, error) {
	if to == u.id {
		return nil, ErrAddYourselfAsFriend
	}

	if u.IsFriend(to) {
		return nil, ErrAlreadyBeenFriends
	}

	return CreateFriendRequest(u.id, to, content, operationIDGenerator), nil
}

func (u *User) ReceiveFriendRequest(
	request *FriendRequest,
	idGenerator event.IDGenerator,
) error {
	if request.to != u.id {
		return ErrFriendRequestNotForUser
	}

	if u.IsFriend(request.from) {
		return ErrAlreadyBeenFriends
	}

	for _, friendRequest := range u.requests {
		if friendRequest.From() == request.From() && friendRequest.state == StatePending {
			return ErrFriendRequestExists
		}
	}

	u.requests = append(u.requests, request)

	ev, err := NewFriendRequestReceivedEvent(
		u.id,
		request.ID(),
		idGenerator,
	)
	if err != nil {
		return err
	}

	u.eventManager.RecordEvent(ev)

	return nil
}

func (u *User) AgreeFriendRequest(
	id kernel.OperationID,
	idGenerator event.IDGenerator,
) error {
	for i, request := range u.requests {
		if request.ID() == id && request.state == StatePending {
			u.requests[i].Agreed()
			u.friendIDs = append(u.friendIDs, request.From())

			ev, err := NewFriendRequestAgreedEvent(
				request.id,
				idGenerator,
			)
			if err != nil {
				return err
			}

			u.eventManager.RecordEvent(ev)

			break
		}
	}
	return nil
}

func (u *User) RefuseFriendRequest(
	id kernel.OperationID,
	idGenerator event.IDGenerator,
) error {
	for i, request := range u.requests {
		if request.ID() == id && request.state == StatePending {
			u.requests[i].Refused()

			ev, err := NewFriendRequestRefusedEvent(
				request.id,
				idGenerator,
			)
			if err != nil {
				return err
			}

			u.eventManager.RecordEvent(ev)

			break
		}
	}
	return nil
}

func (u *User) IsFriend(id kernel.UserID) bool {
	for _, friendID := range u.friendIDs {
		if friendID == id {
			return true
		}
	}
	return false
}

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) Number() kernel.UserNumber {
	return u.number
}

func (u *User) Requests() []*FriendRequest {
	return u.requests
}

func (u *User) FriendIDs() []kernel.UserID {
	return u.friendIDs
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
