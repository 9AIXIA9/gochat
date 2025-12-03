//go:generate mockgen -source=room_repository.go -destination=./mocks/mock_room_repository.go -package=mocks
package domain

type RoomRepository interface {
}
