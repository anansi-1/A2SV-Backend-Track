package router_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "task-manager/Delivery/router"
    "task-manager/Domain"
    "task-manager/Infrastructure"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserController struct {
    mock.Mock
}

func (m *MockUserController) Register(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusCreated, gin.H{"id": "123", "email": "test@example.com"})
}

func (m *MockUserController) Login(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusOK, gin.H{"token": "jwt.token.here"})
}

func (m *MockUserController) PromoteUser(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusOK, gin.H{"id": "123", "role": "admin"})
}

type MockTaskController struct {
    mock.Mock
}

func (m *MockTaskController) GetAllTasks(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusOK, []gin.H{{"id": "1", "title": "Test Task"}})
}

func (m *MockTaskController) GetTaskByID(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusOK, gin.H{"id": "1", "title": "Test Task"})
}

func (m *MockTaskController) CreateTask(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusCreated, gin.H{"id": "1", "title": "New Task"})
}

func (m *MockTaskController) UpdateTask(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusOK, gin.H{"id": "1", "title": "Updated Task"})
}

func (m *MockTaskController) DeleteTask(c *gin.Context) {
    m.Called(c)
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

type UserControllerWrapper struct {
    Mock *MockUserController
}

func (w *UserControllerWrapper) Register(c *gin.Context) {
    w.Mock.Register(c)
}

func (w *UserControllerWrapper) Login(c *gin.Context) {
    w.Mock.Login(c)
}

func (w *UserControllerWrapper) PromoteUser(c *gin.Context) {
    w.Mock.PromoteUser(c)
}

// Wrapper for TaskController
type TaskControllerWrapper struct {
    Mock *MockTaskController
}

func (w *TaskControllerWrapper) GetAllTasks(c *gin.Context) {
    w.Mock.GetAllTasks(c)
}

func (w *TaskControllerWrapper) GetTaskByID(c *gin.Context) {
    w.Mock.GetTaskByID(c)
}

func (w *TaskControllerWrapper) CreateTask(c *gin.Context) {
    w.Mock.CreateTask(c)
}

func (w *TaskControllerWrapper) UpdateTask(c *gin.Context) {
    w.Mock.UpdateTask(c)
}

func (w *TaskControllerWrapper) DeleteTask(c *gin.Context) {
    w.Mock.DeleteTask(c)
}

func setupRouterTest(t *testing.T) (*gin.Engine, *MockUserController, *MockTaskController, func()) {
    t.Helper()

    gin.SetMode(gin.TestMode)

    originalKey := os.Getenv("SECRET_KEY")
    os.Setenv("SECRET_KEY", "test-secret-key-for-router-testing")

    jwtService := Infrastructure.NewJWTService()
    authMiddleware := Infrastructure.NewAuthMiddleware(jwtService)

    mockUserController := new(MockUserController)
    mockTaskController := new(MockTaskController)

    userWrapper := &UserControllerWrapper{Mock: mockUserController}
    taskWrapper := &TaskControllerWrapper{Mock: mockTaskController}

    routerEngine := router.SetupRouter(userWrapper, taskWrapper, authMiddleware)

    cleanup := func() {
        if originalKey != "" {
            os.Setenv("SECRET_KEY", originalKey)
        } else {
            os.Unsetenv("SECRET_KEY")
        }
        gin.SetMode(gin.DebugMode)
    }

    return routerEngine, mockUserController, mockTaskController, cleanup
}

func createValidToken(t *testing.T, email, role string) string {
    t.Helper()
    
    jwtService := Infrastructure.NewJWTService()
    user := Domain.User{Email: email, Role: role}
    token, err := jwtService.GenerateToken(user)
    assert.NoError(t, err)
    return token
}

// user routes tests

func TestRouter_UserRegister_Success(t *testing.T) {
    routerEngine, mockUserController, _, cleanup := setupRouterTest(t)
    defer cleanup()
    
    mockUserController.On("Register", mock.Anything)
    
    reqBody := map[string]string{"email": "test@example.com", "password": "password123"}
    jsonData, _ := json.Marshal(reqBody)
    
    req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    routerEngine.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    mockUserController.AssertCalled(t, "Register", mock.Anything)
}

func TestRouter_UserLogin_Success(t *testing.T) {
    routerEngine, mockUserController, _, cleanup := setupRouterTest(t)
    defer cleanup()
    
    mockUserController.On("Login", mock.Anything)
    
    reqBody := map[string]string{"email": "test@example.com", "password": "password123"}
    jsonData, _ := json.Marshal(reqBody)
    
    req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    routerEngine.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    mockUserController.AssertCalled(t, "Login", mock.Anything)
}

func TestRouter_UserPromote_RequiresAuth(t *testing.T) {
    routerEngine, mockUserController, _, cleanup := setupRouterTest(t)
    defer cleanup()
    

	req := httptest.NewRequest("PUT", "/users/promote/123", nil)
    w := httptest.NewRecorder()
    
    routerEngine.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "missing token", response["message"])
    
    mockUserController.AssertNotCalled(t, "PromoteUser")
}

func TestRouter_UserPromote_RequiresAdmin(t *testing.T) {
    routerEngine, mockUserController, _, cleanup := setupRouterTest(t)
    defer cleanup()
    
    token := createValidToken(t, "user@example.com", "user")
    
    req := httptest.NewRequest("PUT", "/users/promote/123", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w := httptest.NewRecorder()
    
    routerEngine.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusForbidden, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "unauthorized", response["message"])
    
    mockUserController.AssertNotCalled(t, "PromoteUser")
}

