package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTask = "tasks"
	CollectionUser = "users"
	JWTSecret      = "it has to be secret"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	DueDate     time.Time          `json:"due_date"`
	Status      string             `json:"status"`
}
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Role     string             `json:"role"`
}

type TaskRepository interface {
	GetTasks(c context.Context) []Task
	GetTask(c context.Context, taskID primitive.ObjectID) (Task, error)
	AddTask(c context.Context, task *Task) error
	RemoveTask(c context.Context, taskID primitive.ObjectID) error
	UpdateTask(c context.Context, taskID primitive.ObjectID, updatedTask *Task) error
}

type UserRepository interface {
	Register(c context.Context, newUser *User) error
	Login(c context.Context, unauthenticatedUser *User) (string, error)
	Promote(c context.Context, regular_user_ID primitive.ObjectID) error
}

type TaskUsecase interface {
	GetTasks(c context.Context) []Task
	GetTask(c context.Context, taskID primitive.ObjectID) (Task, error)
	AddTask(c context.Context, task *Task) error
	RemoveTask(c context.Context, taskID primitive.ObjectID) error
	UpdateTask(c context.Context, taskID primitive.ObjectID, updatedTask *Task) error
}

type UserUsecase interface {
	Register(c context.Context, newUser *User) error
	Login(c context.Context, unauthenticatedUser *User) (string, error)
	Promote(c context.Context, regular_user_ID primitive.ObjectID) error
}
