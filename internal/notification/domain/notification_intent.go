package contract

type NotificationIntentID string

func (id NotificationIntentID) String() string {
	return string(id)
}

type NotificationIntent struct {
	ID NotificationIntentID
}
