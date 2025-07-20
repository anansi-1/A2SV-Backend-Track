package controllers

import (
	"net/http"

	"github/anansi-1/Task-Four-Task-Manager/models"
	"github/anansi-1/Task-Four-Task-Manager/data"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskService data.TaskService
}

func NewTaskController(taskService data.TaskService) TaskController {
	return TaskController{
		TaskService: taskService,
	}
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := tc.TaskService.AddTask(&task)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func (tc *TaskController) GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := tc.TaskService.GetTask(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (tc *TaskController) GetTasks(ctx *gin.Context) {
	tasks, err := tc.TaskService.GetTasks()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedTask models.Task
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := tc.TaskService.UpdateTask(id, &updatedTask)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	err := tc.TaskService.RemoveTask(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}