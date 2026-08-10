package migrations

import (
	"log"
	"taskflow-api/database"
	"taskflow-api/models"
)

func RunMigrations() {
	err := database.DB.AutoMigrate(
		&models.Task{},
		&models.User{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}
