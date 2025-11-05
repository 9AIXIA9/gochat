package domain

type MessageState string

const (
	MessageStateCreated   MessageState = "created"
	MessageStateDelivered MessageState = "delivered"
	//MessageStateRead      MessageState = "read"
)

func (s MessageState) String() string {
	return string(s)
}
