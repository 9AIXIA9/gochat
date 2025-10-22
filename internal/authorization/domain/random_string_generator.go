package domain

type RandomStringGenerator interface {
	Generate() (string, error)
}
