package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword accepts plain password and returns a hashed string
func HashPassword(plainPassword string, cost int) (string, error) {
	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(plainPassword), cost)
	if err != nil {
		return "", err
	}
	hashedPasswordString := string(hashedPasswordByte)
	return hashedPasswordString, nil
}

// VerifyPassword accepts plain password and hashed password,
// compares them then, returns true if it is a match, false otherwise
func VerifyPassword(plainPassword string, hashedPassword string) bool {
	plainPasswordByte := []byte(plainPassword)
	hashedPasswordByte := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashedPasswordByte, plainPasswordByte)
	return err == nil
}
