package controllers

import (
	"net/http"
	service "task-six-authentication/data"
	"task-six-authentication/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	TaskService service.TaskService
	UserService service.UserService
}

func (controller *Controller) Register(context *gin.Context) {
	var newUser models.User
	if err := context.BindJSON(&newUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		context.Abort()
		return
	}
	err2 := controller.UserService.Register(newUser)
	if err2 != nil {
		context.JSON(500, gin.H{"Error": err2.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (controller *Controller) Login(context *gin.Context) {
	var unauthenticatedUser models.User
	if err := context.BindJSON(&unauthenticatedUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		context.Abort()
		return
	}
	token, err := controller.UserService.Login(unauthenticatedUser)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})

}
func (controller *Controller) Promote(context *gin.Context) {
	regular_user_id := context.Param("_id")
	ObjID, err := primitive.ObjectIDFromHex(regular_user_id)
	if err != nil {
		context.JSON(400, gin.H{"error": "Invalid ID"})
	}
	if err := controller.UserService.Promote(ObjID); err != nil {
		context.JSON(401, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(200, gin.H{"message": "User promoted successfully"})
}
func (controller *Controller) GetTasks(context *gin.Context) {
	tasks := controller.TaskService.GetTasks()
	context.IndentedJSON(http.StatusOK, tasks)
}
func (controller *Controller) GetTask(context *gin.Context) {
	id, err1 := primitive.ObjectIDFromHex(context.Param("_id"))
	if err1 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	task, err2 := controller.TaskService.GetTask(id)
	if err2 != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		return
	}

	context.JSON(http.StatusOK, task)

}
func (controller *Controller) AddTask(context *gin.Context) {
	var newTask models.Task

	if err := context.BindJSON(&newTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.TaskService.AddTask(newTask)
	context.JSON(http.StatusCreated, gin.H{"message": "Task Created"})
}
func (controller *Controller) RemoveTask(context *gin.Context) {
	id, err1 := primitive.ObjectIDFromHex(context.Param("_id"))
	if err1 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err2 := controller.TaskService.RemoveTask(id)
	if err2 == nil {
		context.JSON(http.StatusOK, gin.H{"message": "Task removed"})
		return
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
}
func (controller *Controller) UpdateTask(context *gin.Context) {
	var updatedTask models.Task
	id, err1 := primitive.ObjectIDFromHex(context.Param("_id"))
	if err1 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err2 := context.BindJSON(&updatedTask); err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	orginal_task, err := controller.TaskService.GetTask(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}
	if updatedTask.Title == "" {
		updatedTask.Title = orginal_task.Title
	}
	if updatedTask.Description == "" {
		updatedTask.Description = orginal_task.Description
	}
	if updatedTask.Status == "" {
		updatedTask.Status = orginal_task.Status
	}
	if updatedTask.DueDate.IsZero() {
		updatedTask.DueDate = orginal_task.DueDate
	}
	err2 := controller.TaskService.UpdateTask(id, updatedTask)
	if err2 == nil {
		context.JSON(http.StatusOK, gin.H{"message": "Task Updated"})
		return
	}
	context.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}
