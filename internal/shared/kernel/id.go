package kernel

type ID string

type UserID ID

func (id ID) String() string {
	return string(id)
}

func (id UserID) String() string {
	return string(id)
}
