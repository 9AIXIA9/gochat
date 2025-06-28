package repository

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name         string
	Number       domain.RoomNumber
	Owner        domain.UserNumber
	SecretHash   string
	Description  string
	CurrentUsers int
	MaxUsers     int
}

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *domain.Room) (bool, error) {
	return false, nil
}

func (r *RoomRepository) JoinOne(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return false, nil
}

func (r *RoomRepository) QueryByRoomNumber(roomNumber domain.RoomNumber) (*domain.Room, error) {
	return nil, nil
}

func (r *RoomRepository) Delete(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return false, nil
}
