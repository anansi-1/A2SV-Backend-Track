package Usecases

import (
	"task-manager/Domain"
	"time"
)

type TaskUsecase struct {
	TaskRepo Domain.ITaskRepository
}

func NewTaskUsecase(repo Domain.ITaskRepository) *TaskUsecase {
	return &TaskUsecase{
		TaskRepo: repo,
	}
}

func (u *TaskUsecase) GetAllTasks() ([]Domain.Task, error) {
	return u.TaskRepo.GetAll()
}

func (u *TaskUsecase) GetTaskByID(id string) (Domain.Task, error) {
	return u.TaskRepo.GetByID(id)
}

func (u *TaskUsecase) CreateTask(task Domain.Task) (Domain.Task, error) {
	if task.DueDate.IsZero() {
		task.DueDate = time.Now().UTC()
	}
	return u.TaskRepo.Create(task)
}

func (u *TaskUsecase) UpdateTask(id string, task Domain.Task) (Domain.Task, error) {
	return u.TaskRepo.Update(id, task)
}

func (u *TaskUsecase) DeleteTask(id string) error {
	return u.TaskRepo.Delete(id)
}
