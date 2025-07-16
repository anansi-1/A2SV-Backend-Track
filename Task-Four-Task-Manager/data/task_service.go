package data

import (
	"errors"
	"github/anansi-1/Task-Four-Task-Manager/models"
)

type Task struct {
	tasks []models.Task
}

type TaskService interface {
	GetTasks() []models.Task
	GetTask(Id string) (models.Task, error)
	AddTask(task models.Task)
	UpdateTask(Id string, updatedTask models.Task) string
	RemoveTask(Id string) string
}

func (s *Task) GetTasks() []models.Task {
	return s.tasks
}

func (s *Task) GetTask(id string) (models.Task, error) {
	for _, task := range s.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return models.Task{}, errors.New("no task found")
}

func (s *Task) AddTask(task models.Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Task) RemoveTask(id string) string {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return "Updated"
		}
	}
	return "Failed"
}

func (s *Task) UpdateTask(id string, updatedTask models.Task) string {
	for i, task := range s.tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				task.Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				task.Description = updatedTask.Description
			}
			if !updatedTask.DueDate.IsZero() {
				task.DueDate = updatedTask.DueDate
			}
			if updatedTask.Status != "" {
				task.Status = updatedTask.Status
			}
			s.tasks[i] = task
			return "Updated"
		}
	}
	return "Failed"
}
