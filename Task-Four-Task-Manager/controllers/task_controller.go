package controllers

import (
	"github/anansi-1/Task-Four-Task-Manager/data"
	"net/http"
	"github/anansi-1/Task-Four-Task-Manager/models"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	taskService data.Task
}

func (controller *Controller) GetTasks(context *gin.Context) {
	tasks := controller.taskService.GetTasks()
	context.IndentedJSON(http.StatusOK, tasks)
}

func (controller *Controller) GetTask(context *gin.Context) {
	id := context.Param("id")
	task, err := controller.taskService.GetTask(id)
	if err != nil {
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
	controller.taskService.AddTask(newTask)
	context.JSON(http.StatusCreated, gin.H{"message": "Task Created"})
}

func (controller *Controller) RemoveTask(context *gin.Context) {
	id := context.Param("id")
	result := controller.taskService.RemoveTask(id)
	if result == "Updated" {
		context.JSON(http.StatusOK, gin.H{"message": "Task removed"})
		return
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
}

func (controller *Controller) UpdateTask(context *gin.Context) {
	var updatedTask models.Task
	id := context.Param("id")

	if err := context.BindJSON(&updatedTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := controller.taskService.UpdateTask(id, updatedTask)
	if result == "Updated" {
		context.JSON(http.StatusOK, gin.H{"message": "Task Updated"})
		return
	}
	context.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}
