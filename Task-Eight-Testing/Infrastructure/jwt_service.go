package Infrastructure

import (
	"errors"
	"os"
	"task-manager/Domain"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtCustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type JWTService struct {
	secretKey string
}

func NewJWTService() Domain.IJWTService {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("SECRET_KEY not set in env")
	}
	return &JWTService{secretKey: secretKey}
}

func (j *JWTService) GenerateToken(user Domain.User) (string, error) {
	claims := jwtCustomClaims{
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Subject:   user.Email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenStr string) (*Domain.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &Domain.AuthClaims{
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}
