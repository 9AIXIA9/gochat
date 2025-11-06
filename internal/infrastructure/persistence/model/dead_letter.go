package model

type DeadLetter struct {
	*Event `gorm:"embedded"`
	Reason string `gorm:"type:text;not null"`
}

func NewDeadLetter(event *Event, reason error) *DeadLetter {
	return &DeadLetter{Event: event, Reason: reason.Error()}
}

func (e *DeadLetter) TableName() string {
	return "gochat.dead_letters"
}
