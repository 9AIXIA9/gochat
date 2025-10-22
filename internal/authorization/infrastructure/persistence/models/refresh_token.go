package models

import (
	"gochat/internal/shared/kernel"
	"time"
)

type RefreshToken struct {
	String          string        //unique
	UserID          kernel.UserID //unique
	ExpiredAt       time.Time
	RefreshCount    int
	MaxRefreshCount int
}
