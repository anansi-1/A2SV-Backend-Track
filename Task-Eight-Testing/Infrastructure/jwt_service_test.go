package Infrastructure_test

import (
	"fmt"
	"os"
	"strings"
	"task-manager/Domain"
	"task-manager/Infrastructure"
	"testing"
	

	"github.com/stretchr/testify/assert"
)

func setupJWTTest(t *testing.T) func() {
    t.Helper()
    
    originalKey := os.Getenv("SECRET_KEY")
    
    os.Setenv("SECRET_KEY", "test-secret-key-for-testing-123")
    
    return func() {
        if originalKey != "" {
            os.Setenv("SECRET_KEY", originalKey)
        } else {
            os.Unsetenv("SECRET_KEY")
        }
    }
}

func createTestUser() Domain.User {
    return Domain.User{
        ID:    "123",
        Email: "test@example.com",
        Role:  "user",
    }
}

func TestNewJWTService_Success(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    
    assert.NotNil(t, service, "jwtService should be created successfully")
}

func TestNewJWTService_PanicsWithoutSecretKey(t *testing.T) {
    originalKey := os.Getenv("SECRET_KEY")
    defer func() {
        if originalKey != "" {
            os.Setenv("SECRET_KEY", originalKey)
        } else {
            os.Unsetenv("SECRET_KEY")
        }
    }()
    
    os.Unsetenv("SECRET_KEY")
    
    assert.Panics(t, func() {
        Infrastructure.NewJWTService()
    }, "Should panic when SECRET_KEY is not set")
}

func TestJWTService_GenerateToken_Success(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    user := createTestUser()
    
    token, err := service.GenerateToken(user)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    assert.Contains(t, token, ".", "JWT should contain dots as separators")
    
    parts := len(token) - len(strings.ReplaceAll(token, ".", ""))
    assert.Equal(t, 2, parts, "JWT should have exactly 2 dots (3 parts)")
}

func TestJWTService_GenerateToken_DifferentUsers(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    
    user1 := Domain.User{Email: "user1@example.com", Role: "user"}
    user2 := Domain.User{Email: "user2@example.com", Role: "admin"}
    
    token1, err1 := service.GenerateToken(user1)
    token2, err2 := service.GenerateToken(user2)
    
    assert.NoError(t, err1)
    assert.NoError(t, err2)
    assert.NotEqual(t, token1, token2, "Different users should generate different tokens")
}

func TestJWTService_GenerateToken_EmptyUser(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    emptyUser := Domain.User{} 
    
    token, err := service.GenerateToken(emptyUser)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}

func TestJWTService_ValidateToken_ValidToken(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    user := createTestUser()
    
    token, err := service.GenerateToken(user)
    assert.NoError(t, err)
    
    claims, err := service.ValidateToken(token)
    
    assert.NoError(t, err)
    assert.NotNil(t, claims)
    assert.Equal(t, user.Email, claims.Email)
    assert.Equal(t, user.Role, claims.Role)
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    
    testCases := []struct {
        name  string
        token string
    }{
        {"Completely invalid", "invalid.jwt.token"},
        {"Empty token", ""},
        {"Random string", "just-a-random-string"},
        {"Malformed JWT", "header.payload"},
        {"Too many parts", "header.payload.signature.extra"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            claims, err := service.ValidateToken(tc.token)
            
            assert.Error(t, err)
            assert.Nil(t, claims)
            assert.Contains(t, err.Error(), "invalid", "Error should mention token is invalid")
        })
    }
}

func TestJWTService_ValidateToken_TokenFromDifferentSecret(t *testing.T) {
    cleanup1 := setupJWTTest(t)
    service1 := Infrastructure.NewJWTService()
    user := createTestUser()
    token, _ := service1.GenerateToken(user)
    cleanup1()
    
    os.Setenv("SECRET_KEY", "different-secret-key")
    defer os.Unsetenv("SECRET_KEY")
    
    service2 := Infrastructure.NewJWTService()
    claims, err := service2.ValidateToken(token)
    
    assert.Error(t, err)
    assert.Nil(t, claims)
    assert.Contains(t, err.Error(), "invalid")
}

func TestJWTService_TokenRoundTrip(t *testing.T) {
    cleanup := setupJWTTest(t)
    defer cleanup()
    
    service := Infrastructure.NewJWTService()
    
    testUsers := []Domain.User{
        {Email: "admin@example.com", Role: "admin"},
        {Email: "user@example.com", Role: "user"},
        {Email: "manager@example.com", Role: "manager"},
        {Email: "test@example.com", Role: ""},  
        {Email: "", Role: "user"},              
    }
    
    for i, user := range testUsers {
        t.Run(fmt.Sprintf("User_%d", i), func(t *testing.T) {
            token, err := service.GenerateToken(user)
            assert.NoError(t, err)
            assert.NotEmpty(t, token)
            
            claims, err := service.ValidateToken(token)
            assert.NoError(t, err)
            assert.NotNil(t, claims)
            assert.Equal(t, user.Email, claims.Email)
            assert.Equal(t, user.Role, claims.Role)
        })
    }
}

