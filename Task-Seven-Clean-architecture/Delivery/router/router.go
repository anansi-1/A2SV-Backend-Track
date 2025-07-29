package router

import (
	controller "task-manager/Delivery/controllers"
	"task-manager/Infrastructure"

		"github.com/gin-gonic/gin"

)

func SetupRouter(
    userC *controller.UserController, 
    taskC *controller.TaskController,
    authMiddleware *Infrastructure.AuthMiddleware,
) *gin.Engine {
    router := gin.Default()

    userRoutes := router.Group("/users")
    {
        userRoutes.POST("/register", userC.Register)
        userRoutes.POST("/login", userC.Login)

        userRoutes.PUT("/promote/:id", authMiddleware.Middleware(), authMiddleware.AdminMiddleware(), userC.PromoteUser)
    }

    taskRoutes := router.Group("/tasks", authMiddleware.Middleware())
    {
        taskRoutes.GET("/", taskC.GetAllTasks)
        taskRoutes.GET("/:id", taskC.GetTaskByID)
        taskRoutes.POST("/", authMiddleware.AdminMiddleware(),taskC.CreateTask)
        taskRoutes.PUT("/:id", authMiddleware.AdminMiddleware(),taskC.UpdateTask)
        taskRoutes.DELETE("/:id", authMiddleware.AdminMiddleware(),taskC.DeleteTask)
    }

    return router
}
