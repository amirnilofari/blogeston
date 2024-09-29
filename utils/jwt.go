package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecretKey []byte

func GenerateJWT(userID int, userRole string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["user_role"] = userRole
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	JwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(JwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
