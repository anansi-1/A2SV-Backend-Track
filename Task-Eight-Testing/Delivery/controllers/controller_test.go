package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-manager/Delivery/controllers"
	"task-manager/Domain"
	"task-manager/Usecases"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user Domain.User) (Domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(Domain.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(email string) (Domain.User, error) {
	args := m.Called(email)
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
func (m *MockPasswordService) Compare(raw, hashed string) bool {
	args := m.Called(raw, hashed)
	return args.Bool(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user Domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*Domain.AuthClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*Domain.AuthClaims), args.Error(1)
}


func setupUserController() (*gin.Engine, *MockUserRepository, *MockPasswordService, *MockJWTService) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordService)
	mockJWT := new(MockJWTService)

	usecase := Usecases.NewUserUsecase(mockRepo, mockHasher,mockJWT)
	userController := controllers.NewUserController(usecase)

	r := gin.New()
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.PUT("/users/:id/promote", userController.PromoteUser)

	return r, mockRepo, mockHasher, mockJWT
}


func TestRegister_Success(t *testing.T) {
	r, repo, hasher, _ := setupUserController()
	input := Domain.User{Email: "test@example.com", Password: "pass"}
	hashed := "hashed-pass"
	created := Domain.User{ID: "1", Email: input.Email, Role: "user"}

	repo.On("FindByEmail", input.Email).Return(Domain.User{}, errors.New("not found"))
	hasher.On("Hash", input.Password).Return(hashed, nil)
	repo.On("Create", mock.MatchedBy(func(u Domain.User) bool {
		return u.Email == input.Email && u.Password == hashed
	})).Return(created, nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp Domain.User
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, created.ID, resp.ID)
	assert.Empty(t, resp.Password)
}

func TestLogin_Success(t *testing.T) {
	r, repo, hasher, jwt := setupUserController()
	input := map[string]string{"email": "test@example.com", "password": "pass"}
	user := Domain.User{ID: "1", Email: input["email"], Password: "hashed-pass"}

	repo.On("FindByEmail", input["email"]).Return(user, nil)
	hasher.On("Compare", input["password"], user.Password).Return(true)
	jwt.On("GenerateToken", user).Return("token-abc", nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "token-abc", resp["token"])
}

func TestPromoteUser_Success(t *testing.T) {
	r, repo, _, _ := setupUserController()
	id := "abc123"
	promoted := Domain.User{ID: id, Role: "admin"}
	repo.On("Promote", id).Return(promoted, nil)

	req := httptest.NewRequest(http.MethodPut, "/users/"+id+"/promote", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp Domain.User
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "admin", resp.Role)
}
