package domain

type NotificationState string

const (
	StateUndelivered NotificationState = "undelivered"
	StateDelivered   NotificationState = "delivered"
)
