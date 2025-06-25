package model

import (
	"backend/internal/domain"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name       string
	Number     domain.RoomNumber
	SecretHash string
}
