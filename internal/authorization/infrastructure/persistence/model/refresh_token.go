package model

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RefreshToken struct {
	Token           domain.RefreshToken //unique
	UserID          kernel.UserID       //unique
	ExpiredAt       time.Time
	RefreshCount    int
	MaxRefreshCount int
}
