package routes

import (
	"taskflow-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", handlers.GetTasks)
		taskRoutes.POST("", handlers.CreateTask)
		taskRoutes.GET("/:id", handlers.GetTaskByID)
		taskRoutes.PUT("/:id", handlers.UpdateTask)
		taskRoutes.DELETE("/:id", handlers.DeleteTask)
	}
}
