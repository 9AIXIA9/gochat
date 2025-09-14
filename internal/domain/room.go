package domain

type RoomNumber int64
type Room struct {
	name         string
	number       RoomNumber
	secretHash   string
	description  string
	currentUsers int
	maxUsers     int
	owner        UserNumber
}

func NewRoom(number RoomNumber, name string, secretHash string, description string, currentUsers, maxUsers int, owner UserNumber) *Room {
	return &Room{
		number:       number,
		name:         name,
		secretHash:   secretHash,
		description:  description,
		currentUsers: currentUsers,
		maxUsers:     maxUsers,
		owner:        owner,
	}
}

//Getter

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) SecretHash() string {
	return r.secretHash
}

func (r *Room) Description() string {
	return r.description
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
