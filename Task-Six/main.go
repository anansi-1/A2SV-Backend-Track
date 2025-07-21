package main

import (
	"context"
	"log"
	"task-six-authentication/controllers"
	service "task-six-authentication/data"
	"task-six-authentication/router"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	taskService := service.NewTaskService(client)
	userService := service.NewUserService(client)
	taskController := &controllers.Controller{TaskService: *taskService, UserService: *userService}
	r := router.Router{Controller: taskController}
	r.Route()
}
