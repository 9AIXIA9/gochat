package event

type Topic string

const (
	UserSignedUp Topic = "user.signed.up"
	UserLoggedIn       = "user.logged.in"
)
