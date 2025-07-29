package Infrastructure_test

import (
    "testing"
    "task-manager/Infrastructure"
    
    "github.com/stretchr/testify/assert"
)

func TestPasswordService_Hash_Success(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    password := "mySecretPassword123"
    
    hashed, err := service.Hash(password)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, hashed)
    assert.NotEqual(t, password, hashed, "Hashed password should be different from plain text")
    assert.Greater(t, len(hashed), 50, "Bcrypt hash should be reasonably long")
    
    assert.Contains(t, hashed, "$2a$", "Should contain bcrypt identifier")
}

func TestPasswordService_Hash_EmptyPassword(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    
    hashed, err := service.Hash("")
    
    assert.NoError(t, err)
    assert.NotEmpty(t, hashed)
}


func TestPasswordService_Compare_ValidPassword(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    password := "testPassword123"
    
    hashed, err := service.Hash(password)
    assert.NoError(t, err)
    
    isValid := service.Compare(password, hashed)
    
    assert.True(t, isValid, "Valid password should match its hash")
}

func TestPasswordService_Compare_InvalidPassword(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    correctPassword := "correctPassword"
    wrongPassword := "wrongPassword"
    
    hashed, err := service.Hash(correctPassword)
    assert.NoError(t, err)
    
    isValid := service.Compare(wrongPassword, hashed)
    
    assert.False(t, isValid, "Invalid password should not match hash")
}

func TestPasswordService_Compare_InvalidHash(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    password := "testPassword"
    invalidHash := "not-a-valid-bcrypt-hash"
    
    isValid := service.Compare(password, invalidHash)
    
    assert.False(t, isValid, "Invalid hash should return false")
}

func TestPasswordService_Compare_EmptyPassword(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    
    hashed, _ := service.Hash("")
    
    isValid := service.Compare("", hashed)
    assert.True(t, isValid, "Empty password should match its hash")
    
    nonEmptyHashed, _ := service.Hash("nonempty")
    isValid = service.Compare("", nonEmptyHashed)
    assert.False(t, isValid, "Empty password should not match non-empty hash")
}

func TestPasswordService_HashConsistency(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    password := "consistencyTest"
    
    hash1, err1 := service.Hash(password)
    hash2, err2 := service.Hash(password)
    
    assert.NoError(t, err1)
    assert.NoError(t, err2)
    
    assert.NotEqual(t, hash1, hash2, "Bcrypt should generate different hashes each time due to salt")
    
    assert.True(t, service.Compare(password, hash1), "First hash should validate correctly")
    assert.True(t, service.Compare(password, hash2), "Second hash should validate correctly")
}


func TestPasswordService_CaseSensitivity(t *testing.T) {
    service := Infrastructure.NewPasswordService()
    password := "CaseSensitivePassword"
    
    hashed, _ := service.Hash(password)
    
    testCases := []struct {
        testPassword string
        shouldMatch  bool
    }{
        {"CaseSensitivePassword", true},  
        {"casesensitivepassword", false}, 
        {"CASESENSITIVEPASSWORD", false}, 
        {"CaseSensitivepassword", false}, 
    }
    
    for _, tc := range testCases {
        t.Run("Case_"+tc.testPassword, func(t *testing.T) {
            isValid := service.Compare(tc.testPassword, hashed)
            assert.Equal(t, tc.shouldMatch, isValid, 
                "Password comparison should be case sensitive")
        })
    }
}