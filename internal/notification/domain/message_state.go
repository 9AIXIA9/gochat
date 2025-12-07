package domain

type MessageState string

const (
	MessageStateUndelivered MessageState = "undelivered"
	MessageStateDelivered   MessageState = "delivered"
)

func (s MessageState) String() string {
	return string(s)
}
