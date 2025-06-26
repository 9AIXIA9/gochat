package model

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Number  domain.UserNumber
	Name    string
	PwdHash string
}
