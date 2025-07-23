package infrastructure

import (
	"fmt"
	domain "task-seven/Domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	token.Claims = jwt.MapClaims{
		"_id":      user.ID,
		"username": user.Username,
		"role":     user.Role,
	}

	jwt_token, err := token.SignedString([]byte(domain.JWTSecret))
	if err != nil {
		return "", err
	}
	return jwt_token, nil
}

func CheckToken(authPart string) (*jwt.Token, error) {

	token, err := jwt.Parse(authPart, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(domain.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func FindClaim(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
