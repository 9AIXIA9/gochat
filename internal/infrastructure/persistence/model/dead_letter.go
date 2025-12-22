package model

import "time"

type DeadLetter struct {
	*Event   `gorm:"embedded"`
	Reason   string    `gorm:"type:text;not null"`
	FailedAt time.Time `gorm:"autoCreateTime"`
}

func (e *DeadLetter) TableName() string {
	return "dead_letters"
}
