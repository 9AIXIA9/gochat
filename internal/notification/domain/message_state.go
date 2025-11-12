package domain

type MessageState string

const (
	MessageStateReceived  MessageState = "received"
	MessageStateDelivered MessageState = "delivered"
	MessageStateRead      MessageState = "read"
)

func (s MessageState) String() string {
	return string(s)
}
