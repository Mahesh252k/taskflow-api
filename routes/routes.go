package routes

import (
	"taskflow-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/users/register", handlers.RegisterUser)
	router.POST("/users/login", handlers.LoginUser)
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", handlers.GetTasks)
		taskRoutes.POST("", handlers.CreateTask)
		taskRoutes.GET("/:id", handlers.GetTaskByID)
		taskRoutes.PUT("/:id", handlers.UpdateTask)
		taskRoutes.PATCH("/:id", handlers.UpdateTask)
		taskRoutes.DELETE("/:id", handlers.DeleteTask)
	}
}
