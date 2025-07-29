package Usecases_test

import (
	"errors"
	"testing"
	"time"

	"task-manager/Domain"
	"task-manager/Usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock 
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetAll() ([]Domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByID(id string) (Domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Create(task Domain.Task) (Domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(id string, task Domain.Task) (Domain.Task, error) {
	args := m.Called(id, task)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// tests 

func TestGetAllTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	expected := []Domain.Task{{ID: "1", Title: "Test Task"}}
	mockRepo.On("GetAll").Return(expected, nil)

	tasks, err := usecase.GetAllTasks()

	assert.NoError(t, err)
	assert.Equal(t, expected, tasks)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	taskID := "42"
	expected := Domain.Task{ID: taskID, Title: "Test Task"}
	mockRepo.On("GetByID", taskID).Return(expected, nil)

	task, err := usecase.GetTaskByID(taskID)

	assert.NoError(t, err)
	assert.Equal(t, expected, task)
	mockRepo.AssertExpectations(t)
}


func TestCreateTask_WithExistingDueDate(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	dueDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := Domain.Task{Title: "Pre-scheduled", DueDate: dueDate}
	mockRepo.On("Create", input).Return(input, nil)

	result, err := usecase.CreateTask(input)

	assert.NoError(t, err)
	assert.Equal(t, dueDate, result.DueDate)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	taskID := "123"
	update := Domain.Task{Title: "Updated"}
	mockRepo.On("Update", taskID, update).Return(update, nil)

	result, err := usecase.UpdateTask(taskID, update)

	assert.NoError(t, err)
	assert.Equal(t, update, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	taskID := "123"
	mockRepo.On("Delete", taskID).Return(nil)

	err := usecase.DeleteTask(taskID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskByID_Error(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	mockRepo.On("GetByID", "missing-id").Return(Domain.Task{}, errors.New("not found"))

	_, err := usecase.GetTaskByID("missing-id")

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	mockRepo.AssertExpectations(t)
}
