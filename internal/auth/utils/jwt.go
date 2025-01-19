package utils

import (
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token with claims
func GenerateJWT(jwtSecretKey []byte, userID string, role string) (string, error) {
	// Define the claims
	claims := models.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create JWT with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseJWT parses and validates a JWT token, returning the claims
func ParseJWT(jwtSecretKey []byte, tokenString string) (*models.JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewAuthorizationError("Invalid signing method")
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract and return the claims
	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.NewAuthorizationError("Invalid token")
	}

	return claims, nil
}
