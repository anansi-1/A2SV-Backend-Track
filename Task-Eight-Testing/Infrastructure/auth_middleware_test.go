package Infrastructure_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "task-manager/Domain"
    "task-manager/Infrastructure"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func setupMiddlewareTest(t *testing.T) (*Infrastructure.AuthMiddleware, Domain.IJWTService, func()) {
	t.Helper()
	
    originalKey := os.Getenv("SECRET_KEY")
    os.Setenv("SECRET_KEY", "test-secret-key-for-middleware-testing")
	
    jwtService := Infrastructure.NewJWTService()
    authMiddleware := Infrastructure.NewAuthMiddleware(jwtService)
	
    gin.SetMode(gin.TestMode)
	
    cleanup := func() {
		if originalKey != "" {
			os.Setenv("SECRET_KEY", originalKey)
			} else {
				os.Unsetenv("SECRET_KEY")
			}
			gin.SetMode(gin.DebugMode)
		}
		
		return authMiddleware, jwtService, cleanup
	}
	
	func createTestRouter(middleware gin.HandlerFunc) *gin.Engine {
		router := gin.New()
		router.Use(middleware)
		router.GET("/protected", func(c *gin.Context) {
			email, emailExists := c.Get("email")
			role, roleExists := c.Get("role")
        c.JSON(http.StatusOK, gin.H{
            "message":      "success",
            "email":        email,
            "role":         role,
            "email_exists": emailExists,
            "role_exists":  roleExists,
        })
    })
    return router
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	authMiddleware, jwtService, cleanup := setupMiddlewareTest(t)
    defer cleanup()
	
    user := Domain.User{Email: "test@example.com", Role: "user"}
    token, err := jwtService.GenerateToken(user)
    assert.NoError(t, err)
	
    router := createTestRouter(authMiddleware.Middleware())
	
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w := httptest.NewRecorder()
	
    router.ServeHTTP(w, req)
	
    assert.Equal(t, http.StatusOK, w.Code)
	
    var response map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
	
    assert.Equal(t, "success", response["message"])
    assert.Equal(t, user.Email, response["email"])
    assert.Equal(t, user.Role, response["role"])
    assert.True(t, response["email_exists"].(bool))
    assert.True(t, response["role_exists"].(bool))
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	authMiddleware, _, cleanup := setupMiddlewareTest(t)
    defer cleanup()
	
    router := createTestRouter(authMiddleware.Middleware())
	
    req := httptest.NewRequest("GET", "/protected", nil)
    w := httptest.NewRecorder()
	
    router.ServeHTTP(w, req)
	
    assert.Equal(t, http.StatusUnauthorized, w.Code)
	
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "missing token", response["message"])
}



func TestAuthMiddleware_InvalidToken(t *testing.T) {
	authMiddleware, _, cleanup := setupMiddlewareTest(t)
    defer cleanup()
	
    router := createTestRouter(authMiddleware.Middleware())
	
    invalidTokens := []string{
		"invalid.jwt.token",
        "completely-invalid-token",
        "",
        "header.payload",
    }
	
    for _, invalidToken := range invalidTokens {
		limit := min(10, len(invalidToken))
        t.Run("Token_"+invalidToken[:limit], func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
            req.Header.Set("Authorization", "Bearer "+invalidToken)
            w := httptest.NewRecorder()
			
            router.ServeHTTP(w, req)
			
            assert.Equal(t, http.StatusForbidden, w.Code)
			
            var response map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &response)
            assert.Equal(t, "unauthorized", response["message"])
        })
    }
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}