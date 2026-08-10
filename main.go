package main

import (
	"taskflow-api/config"
	"taskflow-api/database"
	"taskflow-api/middleware"
	"taskflow-api/migrations"
	"taskflow-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.ConnectDatabase()
	migrations.RunMigrations()
	migrations.Seed()

	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.Auth())

	routes.SetupRoutes(router)

	if err := router.Run(":" + config.GetEnv("APP_PORT")); err != nil {
		panic(err)
	}
}
