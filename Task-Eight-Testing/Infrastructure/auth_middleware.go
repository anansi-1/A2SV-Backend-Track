package Infrastructure

import (
	"net/http"
	"strings"
	"task-manager/Domain"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService Domain.IJWTService
}

func NewAuthMiddleware(jwtService Domain.IJWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (a *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
			return
		}

		authParts := strings.Split(auth, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
			return
		}

		claims, err := a.jwtService.ValidateToken(authParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "unauthorized"})
			return
		}

		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (a *AuthMiddleware) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "unauthorized"})
			return
		}
		c.Next()
	}
}
