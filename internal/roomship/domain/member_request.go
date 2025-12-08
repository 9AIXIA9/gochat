package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const maxContentLength = 100

const (
	StatePending MemberRequestState = "pending"
	StateAgreed  MemberRequestState = "agreed"
	StateRefused MemberRequestState = "refused"
)

type MemberRequestState string

func (s MemberRequestState) String() string {
	return string(s)
}

type MemberRequest struct {
	id          kernel.OperationID
	state       MemberRequestState
	applicantID kernel.UserID
	roomID      kernel.RoomID
	content     string
	operatorID  kernel.UserID
	operatedAt  time.Time //UTC
	createdAt   time.Time //UTC

	manager *event.Manager
}

func LoadMemberRequest(
	id kernel.OperationID,
	state MemberRequestState,
	applicantID kernel.UserID,
	roomID kernel.RoomID,
	content string,
	operatorID kernel.UserID,
	operatedAt time.Time,
	createdAt time.Time,
) *MemberRequest {
	return &MemberRequest{
		id:          id,
		state:       state,
		applicantID: applicantID,
		roomID:      roomID,
		content:     content,
		operatorID:  operatorID,
		operatedAt:  operatedAt,
		createdAt:   createdAt,
		manager:     event.NewEventManager(),
	}
}

func CreateMemberRequest(
	applicantID kernel.UserID,
	roomID kernel.RoomID,
	content string,
	operationIDGenerator kernel.OperationIDGenerator,
	idGenerator event.IDGenerator,
) (*MemberRequest, error) {
	if len(content) > maxContentLength {
		return nil, ErrContentTooLong
	}

	request := &MemberRequest{
		id:          operationIDGenerator.Generate(),
		state:       StatePending,
		applicantID: applicantID,
		roomID:      roomID,
		content:     content,
		operatorID:  "",
		operatedAt:  time.Now().UTC(),
		createdAt:   time.Now().UTC(),
		manager:     event.NewEventManager(),
	}

	ev, err := NewMemberRequestCreatedEvent(request.id, idGenerator)
	if err != nil {
		return nil, err
	}

	request.manager.RecordEvent(ev)

	return request, nil
}

func (r *MemberRequest) Agree(
	operatorID kernel.UserID,
	idGenerator event.IDGenerator,
) error {
	if r.state != StatePending {
		return ErrHandleNotPendingRequest
	}

	ev, err := NewMemberRequestAgreedEvent(r.id, idGenerator)
	if err != nil {
		return err
	}

	r.state = StateAgreed
	r.operatorID = operatorID
	r.operatedAt = time.Now().UTC().
		UTC()

	r.manager.RecordEvent(ev)

	return nil
}

func (r *MemberRequest) Refuse(
	operatorID kernel.UserID,
) error {
	if r.state != StatePending {
		return ErrHandleNotPendingRequest
	}

	r.state = StateRefused
	r.operatorID = operatorID
	r.operatedAt = time.Now().UTC().
		UTC()

	return nil
}

func (r *MemberRequest) ID() kernel.OperationID {
	return r.id
}

func (r *MemberRequest) State() MemberRequestState {
	return r.state
}

func (r *MemberRequest) ApplicantID() kernel.UserID {
	return r.applicantID
}

func (r *MemberRequest) RoomID() kernel.RoomID {
	return r.roomID
}

func (r *MemberRequest) Content() string {
	return r.content
}

func (r *MemberRequest) OperatorID() kernel.UserID {
	return r.operatorID
}

func (r *MemberRequest) OperatedAt() time.Time {
	return r.operatedAt
}

func (r *MemberRequest) CreatedAt() time.Time {
	return r.createdAt
}

func (r *MemberRequest) GetEvents() []event.Event {
	return r.manager.GetEvents()
}
