package usecases

import (
	"context"
	domain "task-seven/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (tu *taskUsecase) GetTasks(c context.Context) []domain.Task {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetTasks(ctx)
}

func (tu *taskUsecase) GetTask(c context.Context, taskID primitive.ObjectID) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetTask(ctx, taskID)
}

func (tu *taskUsecase) AddTask(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.AddTask(ctx, task)
}

func (tu *taskUsecase) RemoveTask(c context.Context, taskID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.RemoveTask(ctx, taskID)
}

func (tu *taskUsecase) UpdateTask(c context.Context, taskID primitive.ObjectID, updatedTask *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.UpdateTask(ctx, taskID, updatedTask)
}
