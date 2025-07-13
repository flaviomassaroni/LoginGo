package main

import (
	"fmt"
	"regexp"

	"time"

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

func UsernameExists(username string, db *sqlx.DB) bool {

	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username)

	if err != nil {
		return false
	}

	return exists

}

func CreatingJWTToken(username string, secretKey []byte) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
