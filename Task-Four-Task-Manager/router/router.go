package router

import (
	"github.com/gin-gonic/gin"
	"github/anansi-1/Task-Four-Task-Manager/controllers"

)

type Router struct {
	controller controllers.Controller
}


func (r *Router) StartRoute() {
	router := gin.Default()
	router.GET("/tasks", r.controller.GetTasks)
	router.GET("/tasks/:id", r.controller.GetTask)
	router.POST("/tasks", r.controller.AddTask)
	router.PUT("/tasks/:id", r.controller.UpdateTask)
	router.DELETE("/tasks/:id", r.controller.RemoveTask)
	router.Run(":8080")
}
