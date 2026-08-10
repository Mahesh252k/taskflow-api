package migrations

import (
	"log"

	"taskflow-api/database"
	"taskflow-api/models"
)

func Seed() {
	user := models.User{
		Name:  "Mike",
		Email: "mike@example.com",
	}

	if err := database.DB.
		Where("email = ?", user.Email).
		FirstOrCreate(&user).Error; err != nil {
		log.Printf("failed to seed user: %v", err)
		return
	}

	var taskCount int64

	if err := database.DB.
		Model(&models.Task{}).
		Where("user_id = ?", user.ID).
		Count(&taskCount).Error; err != nil {
		log.Printf("failed to count tasks: %v", err)
		return
	}

	if taskCount > 0 {
		return
	}

	tasks := []models.Task{
		{
			Title:       "Learn Go",
			Description: "Learn Go fundamentals",
			UserID:      user.ID,
			Done:        false,
		},
		{
			Title:       "Learn GORM",
			Description: "Learn database operations using GORM",
			UserID:      user.ID,
			Done:        false,
		},
	}

	if err := database.DB.Create(&tasks).Error; err != nil {
		log.Printf("failed to seed tasks: %v", err)
	}
}
