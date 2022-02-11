package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//HashPassword trasforma la pwd in un hash usando bcrypt
func HashPassword(password string) (string, error) {
	bpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("Errore in fase di creazione dell'hash della password (%v)", err)
	}
	return string(bpassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
