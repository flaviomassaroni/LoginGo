package main

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckUsername(username string) bool {

	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if len(username) < 3 || len(username) > 16 || !re.MatchString(username) {
		return false
	}

	return true
}

func CheckPassword(password string) bool {

	len := len(password) >= 8

	upper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	lower := regexp.MustCompile(`[a-z]`).MatchString(password)
	nums := regexp.MustCompile(`[0-9]`).MatchString(password)
	special := regexp.MustCompile(`[!@#\$%\^&\*_\-+=\[\]{};':"\\|,.<>/?]`).MatchString(password)

	if !len || !upper || !lower || !nums || !special {
		return false
	}

	return true
}
