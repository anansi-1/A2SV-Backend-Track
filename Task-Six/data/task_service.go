package data

	import ("github/anansi-1/Task-Four-Task-Manager/models")

type TaskService interface {
	GetTasks() ([]*models.Task, error)
	GetTask(id string) (*models.Task, error)
	AddTask(task *models.Task) error
	UpdateTask(id string, updatedTask *models.Task) error
	RemoveTask(id string) error
}