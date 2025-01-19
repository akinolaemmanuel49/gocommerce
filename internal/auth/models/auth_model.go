package models

import "github.com/golang-jwt/jwt/v5"

type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	Token string `json:"token"`
}

type JWTClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
