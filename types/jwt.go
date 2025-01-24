package types

import "github.com/golang-jwt/jwt"

type JWTType string

var (
	JWTTypeRefresh JWTType = "refresh"
	JWTTypeAccess  JWTType = "access"
)

type Claims struct {
	Type   JWTType `json:"type"`
	UserID string  `json:"user_id"`
	jwt.StandardClaims
}
