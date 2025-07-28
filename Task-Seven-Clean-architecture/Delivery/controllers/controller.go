package controllers

import (
	"net/http"
	"task-manager/Domain"
	"task-manager/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)


type UserController struct {
	UserUsecase *Usecases.UserUsecase
}

func NewUserController(userUsecase *Usecases.UserUsecase) *UserController {
	return &UserController{UserUsecase: userUsecase}
}

func (c *UserController) Register(ctx *gin.Context) {
	var user Domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdUser, err := c.UserUsecase.Register(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser.Password = ""
	ctx.JSON(http.StatusCreated, createdUser)
}

func (c *UserController) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password required"})
		return
	}

	token, err := c.UserUsecase.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) PromoteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	updatedUser, err := c.UserUsecase.PromoteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}

	updatedUser.Password = ""
	ctx.JSON(http.StatusOK, updatedUser)
}


type TaskController struct {
	TaskUsecase *Usecases.TaskUsecase
}

func NewTaskController(taskUsecase *Usecases.TaskUsecase) *TaskController {
	return &TaskController{TaskUsecase: taskUsecase}
}

func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	tasks, err := c.TaskUsecase.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.TaskUsecase.GetTaskByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task Domain.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if task.DueDate.IsZero() {
		task.DueDate = time.Now().UTC()
	}

	createdTask, err := c.TaskUsecase.CreateTask(task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, createdTask)
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var task Domain.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedTask, err := c.TaskUsecase.UpdateTask(id, task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	ctx.JSON(http.StatusOK, updatedTask)
}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.TaskUsecase.DeleteTask(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
