package event

type IDGenerator interface {
	Generate() ID
}
