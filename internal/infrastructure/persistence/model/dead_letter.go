package model

type DeadLetter struct {
	*Event `gorm:"embedded"`
	Reason string `gorm:"type:text;not null"`
}

func (e *DeadLetter) TableName() string {
	return "dead_letters"
}
