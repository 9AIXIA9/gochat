//go:generate mockgen -source=room_profile_repository.go -destination=./mocks/mock_room_profile_repository.go -package=mocks
package domain

type RoomProfileRepository interface {
}
