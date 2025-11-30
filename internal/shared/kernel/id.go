package kernel

const (
	EmptyUserID UserID = ""
)

type ID string

func (id ID) String() string {
	return string(id)
}

type UserID ID

func (id UserID) String() string {
	return string(id)
}

type RoomID ID

func (i RoomID) String() string {
	return string(i)
}

type MessageID ID

func (i MessageID) String() string {
	return string(i)
}

type OperationID ID

func (i OperationID) String() string {
	return string(i)
}
