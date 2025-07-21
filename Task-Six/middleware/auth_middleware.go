package middleware

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	jwt_secret := []byte("It has to be a secret")
	return func(context *gin.Context) {
		defer context.Next()

		authHeader := context.GetHeader("Authorization")

		if authHeader == "" {
			context.JSON(401, gin.H{"Error": "Authorization header is required"})
			context.Abort()
			return
		}

		authPart := strings.Split(authHeader, " ")

		if len(authPart) != 2 || strings.ToLower(authPart[0]) != "bearer" {
			context.JSON(401, gin.H{"message": "Invalid Authorization header"})
			context.Abort()
			return
		}

		token, err := jwt.Parse(authPart[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwt_secret, nil
		})

		if err != nil || !token.Valid {
			context.JSON(401, gin.H{"error": "Invalid token"})
			context.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			context.JSON(401, gin.H{"error": "Invalid token claims"})
			context.Abort()
			return
		}

		role, ok := claims["role"]
		if !ok {
			context.JSON(401, gin.H{"error": "Invalid role field in token"})
			context.Abort()
			return
		}

		context.Set("role", role)
	}
}
func AdminMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer context.Next()
		role, exists := context.Get("role")
		if !exists || role != "admin" {
			context.JSON(403, gin.H{"message": "Sorry, you are not eligible to do this.", "your role is": role})
			context.Abort()
			return
		}
	}
}
