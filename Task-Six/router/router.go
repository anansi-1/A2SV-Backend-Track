package router

import (
	controller "task-six-authentication/controllers"
	"task-six-authentication/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller *controller.Controller
}

func (r *Router) Route() {
	router := gin.Default()
	router.POST("/register", r.Controller.Register)
	router.POST("/login", r.Controller.Login)
	router.GET("/tasks", middleware.AuthMiddleWare(), r.Controller.GetTasks)
	router.GET("/tasks/:_id", middleware.AuthMiddleWare(), r.Controller.GetTask)
	router.POST("users/:_id", middleware.AuthMiddleWare(), middleware.AdminMiddleWare(), r.Controller.Promote)
	router.POST("/tasks", middleware.AuthMiddleWare(), middleware.AdminMiddleWare(), r.Controller.AddTask)
	router.PUT("/tasks/:_id", middleware.AuthMiddleWare(), middleware.AdminMiddleWare(), r.Controller.UpdateTask)
	router.DELETE("/tasks/:_id", middleware.AuthMiddleWare(), middleware.AdminMiddleWare(), r.Controller.RemoveTask)
	router.Run(":8080")
}
