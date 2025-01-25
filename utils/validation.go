package utils

import (
	"net/mail"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
