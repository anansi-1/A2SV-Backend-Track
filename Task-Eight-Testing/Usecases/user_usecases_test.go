package Usecases_test

import (
	"errors"
	"testing"
	"task-manager/Domain"
	"task-manager/Usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (Domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(user Domain.User) (Domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) Promote(id string) (Domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(Domain.User), args.Error(1)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) Compare(plain, hashed string) bool {
	args := m.Called(plain, hashed)
	return args.Bool(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user Domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*Domain.AuthClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*Domain.AuthClaims), args.Error(1)
}

// test

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	user := Domain.User{Email: "anansi@test.com", Password: "password123"}
	hashedPassword := "hashed_password"
	expectedUser := Domain.User{ID: "1", Email: "anansi@test.com", Password: hashedPassword}

	mockRepo.On("FindByEmail", user.Email).Return(Domain.User{}, errors.New("not found"))
	mockHasher.On("Hash", user.Password).Return(hashedPassword, nil)
	mockRepo.On("Create", mock.MatchedBy(func(u Domain.User) bool {
		return u.Email == user.Email && u.Password == hashedPassword
	})).Return(expectedUser, nil)

	result, err := usecase.Register(user)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	user := Domain.User{Email: "existing@example.com", Password: "password123"}
	existingUser := Domain.User{ID: "1", Email: "existing@example.com"}

	mockRepo.On("FindByEmail", user.Email).Return(existingUser, nil)

	result, err := usecase.Register(user)

	assert.Error(t, err)
	assert.EqualError(t, err, "email already registered")
	assert.Equal(t, Domain.User{}, result)
	mockRepo.AssertExpectations(t)
}

func TestRegister_HashPasswordFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	user := Domain.User{Email: "new@example.com", Password: "password123"}

	mockRepo.On("FindByEmail", user.Email).Return(Domain.User{}, errors.New("not found"))
	mockHasher.On("Hash", user.Password).Return("", errors.New("hash failure"))

	result, err := usecase.Register(user)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to hash password")
	assert.Equal(t, Domain.User{}, result)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestRegister_CreateUserFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	user := Domain.User{Email: "new@example.com", Password: "password123"}
	hashedPassword := "hashed_password"

	mockRepo.On("FindByEmail", user.Email).Return(Domain.User{}, errors.New("not found"))
	mockHasher.On("Hash", user.Password).Return(hashedPassword, nil)
	mockRepo.On("Create", mock.Anything).Return(Domain.User{}, errors.New("db error"))

	result, err := usecase.Register(user)

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
	assert.Equal(t, Domain.User{}, result)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	email := "anansi@test.com"
	password := "password123"
	hashedPassword := "hashed_password"
	user := Domain.User{ID: "1", Email: email, Password: hashedPassword}
	expectedToken := "jwt_token"

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockHasher.On("Compare", password, hashedPassword).Return(true)
	mockJWT.On("GenerateToken", user).Return(expectedToken, nil)

	token, err := usecase.Login(email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	email := "anansi@test.com"
	password := "wrong_password"
	hashedPassword := "hashed_password"
	user := Domain.User{ID: "1", Email: email, Password: hashedPassword}

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockHasher.On("Compare", password, hashedPassword).Return(false)

	token, err := usecase.Login(email, password)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid credentials")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	email := "missing@example.com"

	mockRepo.On("FindByEmail", email).Return(Domain.User{}, errors.New("not found"))

	token, err := usecase.Login(email, "any_password")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid credentials")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_TokenGenerationFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	email := "anansi@test.com"
	password := "password123"
	hashedPassword := "hashed_password"
	user := Domain.User{ID: "1", Email: email, Password: hashedPassword}

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockHasher.On("Compare", password, hashedPassword).Return(true)
	mockJWT.On("GenerateToken", user).Return("", errors.New("token error"))

	token, err := usecase.Login(email, password)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to generate token")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestPromoteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	userID := "123"
	promotedUser := Domain.User{ID: userID, Email: "promoted@example.com", Role: "admin"}

	mockRepo.On("Promote", userID).Return(promotedUser, nil)

	result, err := usecase.PromoteUser(userID)

	assert.NoError(t, err)
	assert.Equal(t, promotedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestPromoteUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)
	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher, mockJWT)

	userID := "123"

	mockRepo.On("Promote", userID).Return(Domain.User{}, errors.New("promotion failed"))

	result, err := usecase.PromoteUser(userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "promotion failed")
	assert.Equal(t, Domain.User{}, result)
	mockRepo.AssertExpectations(t)
}
