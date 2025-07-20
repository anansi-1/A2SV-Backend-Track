package router

import (
	"context"
	"github/anansi-1/Task-Four-Task-Manager/controllers"
	"github/anansi-1/Task-Four-Task-Manager/data"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Router struct {
	TaskController controllers.TaskController
}

func NewRouter(tc controllers.TaskController) *Router {
	return &Router{
		TaskController: tc,
	}
}

func (r *Router) Start() error {
	router := gin.Default()

	router.GET("/tasks", r.TaskController.GetTasks)
	router.GET("/tasks/:id", r.TaskController.GetTask)
	router.POST("/tasks", r.TaskController.CreateTask)
	router.PUT("/tasks/:id", r.TaskController.UpdateTask)
	router.DELETE("/tasks/:id", r.TaskController.DeleteTask)

	return router.Run(":8080")
}


func Run() error {
	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	defer client.Disconnect(ctx)

	taskCollection := client.Database("taskdb").Collection("tasks")
	taskService := data.NewTaskService(taskCollection, ctx)
	taskController := controllers.NewTaskController(taskService)
	r := NewRouter(taskController)

	return r.Start()
}