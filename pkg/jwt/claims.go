package jwt

import jwt2 "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt2.RegisteredClaims
	UserID   string
	Role     string
	Sections []string // accessible sections by user role
}
