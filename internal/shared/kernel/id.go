package kernel

type ID string

func (id ID) String() string {
	return string(id)
}

type UserID string

func (id UserID) String() string {
	return string(id)
}

type RoomID string

func (i RoomID) String() string {
	return string(i)
}

type MessageID string

func (i MessageID) String() string {
	return string(i)
}

type OperationID string

func (i OperationID) String() string {
	return string(i)
}
