package domain

type RoomNumber BaseNumber
type Room struct {
	number       RoomNumber
	owner        UserNumber
	secretHash   string
	currentUsers int
	maxUsers     int
}

func NewRoom(number RoomNumber, secretHash string, currentUsers, maxUsers int, owner UserNumber) *Room {
	return &Room{
		number:       number,
		secretHash:   secretHash,
		currentUsers: currentUsers,
		maxUsers:     maxUsers,
		owner:        owner,
	}
}

func CreateRoom(number BaseNumber, secretHash string, maxUsers int, owner UserNumber) *Room {
	return NewRoom(RoomNumber(number), secretHash, 0, maxUsers, owner)
}

//Getter

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) SecretHash() string {
	return r.secretHash
}

func (r *Room) CurrentUsers() int {
	return r.currentUsers
}

func (r *Room) MaxUsers() int {
	return r.maxUsers
}

func (r *Room) Owner() UserNumber {
	return r.owner
}
