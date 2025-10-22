package kernel

type ID string

type UserID ID

type IDGenerator interface {
	Generate() ID
}
