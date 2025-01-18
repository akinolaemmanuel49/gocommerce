package models

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	AccessToken  string
	RefreshToken string
}

type JWTClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
