package main

import (
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jmoiron/sqlx"
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

func UsernameExists(w http.ResponseWriter, username string, db *sqlx.DB) bool {

	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	return exists

}

func CreatingJWTToken(username string) (string, error) {

	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	return token
}
