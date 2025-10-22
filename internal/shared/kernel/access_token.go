package kernel

type AccessToken string

type AccessTokenParser interface {
	Parse(tokenString AccessToken) (UserID, error)
}

type AccessTokenGenerator interface {
	Generate(userID UserID) (AccessToken, error)
}
