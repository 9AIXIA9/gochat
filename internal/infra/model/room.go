package model

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name       string
	Number     domain.RoomNumber
	SecretHash string
}
