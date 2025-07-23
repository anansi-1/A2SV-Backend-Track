package controller

import (
	"net/http"
	domain "task-seven/Domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	TaskUsecase domain.TaskUsecase
	UserUsecase domain.UserUsecase
}

func (controller *Controller) Register(c *gin.Context) {
	var newUser *domain.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		c.Abort()
		return
	}
	err2 := controller.UserUsecase.Register(c, newUser)
	if err2 != nil {
		c.JSON(500, gin.H{"Error": err2.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (controller *Controller) Login(c *gin.Context) {
	var unauthenticatedUser *domain.User
	if err := c.BindJSON(&unauthenticatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		c.Abort()
		return
	}
	token, err := controller.UserUsecase.Login(c, unauthenticatedUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})

}
func (controller *Controller) Promote(c *gin.Context) {
	regular_user_id := c.Param("_id")
	ObjID, err := primitive.ObjectIDFromHex(regular_user_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
	}
	if err := controller.UserUsecase.Promote(c, ObjID); err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message": "User promoted successfully"})
}
func (controller *Controller) GetTasks(c *gin.Context) {
	tasks := controller.TaskUsecase.GetTasks(c)
	c.IndentedJSON(http.StatusOK, tasks)
}
func (controller *Controller) GetTask(c *gin.Context) {
	id, err1 := primitive.ObjectIDFromHex(c.Param("_id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	task, err2 := controller.TaskUsecase.GetTask(c, id)
	if err2 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		return
	}

	c.JSON(http.StatusOK, task)

}
func (controller *Controller) AddTask(c *gin.Context) {
	var newTask *domain.Task

	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.TaskUsecase.AddTask(c, newTask)
	c.JSON(http.StatusCreated, gin.H{"message": "Task Created"})
}
func (controller *Controller) RemoveTask(c *gin.Context) {
	id, err1 := primitive.ObjectIDFromHex(c.Param("_id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err2 := controller.TaskUsecase.RemoveTask(c, id)
	if err2 == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
}
func (controller *Controller) UpdateTask(c *gin.Context) {
	var updatedTask *domain.Task
	id, err1 := primitive.ObjectIDFromHex(c.Param("_id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err2 := c.BindJSON(&updatedTask); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	orginal_task, err := controller.TaskUsecase.GetTask(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
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
	err2 := controller.TaskUsecase.UpdateTask(c, id, updatedTask)
	if err2 == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Task Updated"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}
