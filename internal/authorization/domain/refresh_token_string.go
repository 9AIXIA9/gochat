package domain

type RefreshTokenString string

type RefreshTokenGenerator interface {
	Generate() (RefreshTokenString, error)
}
