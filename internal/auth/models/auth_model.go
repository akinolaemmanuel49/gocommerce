package models

import "github.com/golang-jwt/jwt/v5"

type UserCredentials struct {
	Email    string `json:"email" example:"john.doe@example.com" validate:"required,email"`
	Password string `json:"password" example:"password" validate:"required"`
}

type Token struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI2NzhkMTQzNDU2ZDAyYzFjYzI0ODMwODkiLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3Mzc1NDQ0MTMsImlhdCI6MTczNzQ1ODAxM30.-g29ffyyjSkV5oB8RXzq-aydW78LBETLGdCPQoOjjH4"`
}

type JWTClaims struct {
	UserID string `json:"userId" example:"678d143456d02c1cc2483089"`
	Role   string `json:"role" example:"admin"`
	jwt.RegisteredClaims
}
