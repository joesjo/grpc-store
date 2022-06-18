package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims["username"].(string), nil
	}
	return "", err
}
