package domain

type SignupResponse struct {
	UserNumber UserNumber `json:"user_number"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateRoomResponse struct {
	RoomNumber  RoomNumber `json:"room_number,string"`
	RoomName    string     `json:"room_name"`
	Owner       UserNumber `json:"owner,string"`
	Description string     `json:"description"`
	MaxUsers    int        `json:"max_users"`
}
