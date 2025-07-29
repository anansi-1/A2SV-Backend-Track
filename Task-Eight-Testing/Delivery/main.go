package main

import (
	"log"
	"task-manager/Delivery/controllers"
	"task-manager/Delivery/router"
	"task-manager/Infrastructure"
	"task-manager/Repositories"
	"task-manager/Usecases"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    passwordService := Infrastructure.NewPasswordService()
    jwtService := Infrastructure.NewJWTService()

    authMiddleware := Infrastructure.NewAuthMiddleware(jwtService) 

    userRepo := Repositories.NewUserRepository()
    taskRepo := Repositories.NewTaskRepository()

    userUC := Usecases.NewUserUsecase(userRepo, passwordService, jwtService)
    taskUC := Usecases.NewTaskUsecase(taskRepo)

    userController := controllers.NewUserController(userUC)
    taskController := controllers.NewTaskController(taskUC)

    r := router.SetupRouter(userController, taskController, authMiddleware) 
    r.Run()
}
