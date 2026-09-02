package routes

import (
	"taskflow-api/handlers"
	"taskflow-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/users/register", handlers.RegisterUser)
	router.POST("/users/login", handlers.LoginUser)
	adminRoutes := router.Group("/admin")
	{
		adminRoutes.Use(middleware.Auth())
		adminRoutes.Use(middleware.AdminOnly())

		adminRoutes.GET("/tasks", handlers.GetAllTasks)
	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.Use(middleware.Auth())

		taskRoutes.GET("", handlers.GetTasks)
		taskRoutes.POST("", handlers.CreateTask)
		taskRoutes.GET("/:id", handlers.GetTaskByID)
		taskRoutes.PUT("/:id", handlers.UpdateTask)
		taskRoutes.PATCH("/:id", handlers.UpdateTask)
		taskRoutes.DELETE("/:id", handlers.DeleteTask)
	}
}
