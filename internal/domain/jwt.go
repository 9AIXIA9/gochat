package domain

import "github.com/golang-jwt/jwt/v4"

const UserIDKey = "userID"
const UserNumberKey = "userNumber"

type JwtCustomClaims struct {
	UserID     string
	UserNumber UserNumber
	jwt.RegisteredClaims
}
