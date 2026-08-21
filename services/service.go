package services

import (
	"errors"
	"strings"

	"taskflow-api/database"
	"taskflow-api/models"

	"gorm.io/gorm"
)

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrTitleRequired = errors.New("title is required")
	ErrUserRequired  = errors.New("user ID is required")
)

func GetTasks(page, limit int, filter models.TaskFilter) (models.TaskListResponse, error) {
	var tasks []models.Task
	var total int64

	query := database.DB.Model(&models.Task{}).
		Preload("User")

	if filter.Done != nil {
		query = query.Where("done = ?", *filter.Done)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"

		query = query.Where(
			"(title LIKE ? OR description LIKE ?)",
			search,
			search,
		)
	}

	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if err := query.Count(&total).Error; err != nil {
		return models.TaskListResponse{}, err
	}

	query = query.Order(filter.Sort + " " + filter.Order)

	query = query.
		Offset((page - 1) * limit).
		Limit(limit)

	if err := query.Find(&tasks).Error; err != nil {
		return models.TaskListResponse{}, err
	}

	pagination := models.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	return models.TaskListResponse{
		Data:       tasks,
		Pagination: pagination,
	}, nil
}

func CreateTask(task models.Task) (models.Task, error) {
	task.Title = strings.TrimSpace(task.Title)
	task.Description = strings.TrimSpace(task.Description)

	if task.Title == "" {
		return models.Task{}, ErrTitleRequired
	}

	if task.UserID == 0 {
		return models.Task{}, ErrUserRequired
	}

	var user models.User

	err := database.DB.First(&user, task.UserID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, ErrUserNotFound
	}

	if err != nil {
		return models.Task{}, err
	}

	if err := database.DB.Create(&task).Error; err != nil {
		return models.Task{}, err
	}

	task.User = user

	return task, nil
}

func GetTaskByID(id uint) (models.Task, error) {
	var task models.Task

	err := database.DB.
		Preload("User").
		First(&task, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, ErrTaskNotFound
	}

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func UpdateTask(id uint, request models.UpdateTaskRequest) (models.TaskResponse, error) {
	var existingTask models.Task

	err := database.DB.First(&existingTask, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.TaskResponse{}, ErrTaskNotFound
	}

	if err != nil {
		return models.TaskResponse{}, err
	}

	updates := map[string]interface{}{}

	if request.Title != nil {
		title := strings.TrimSpace(*request.Title)

		if title == "" {
			return models.TaskResponse{}, ErrTitleRequired
		}

		updates["title"] = title
	}

	if request.Description != nil {
		updates["description"] = strings.TrimSpace(*request.Description)
	}

	if request.Done != nil {
		updates["done"] = *request.Done
	}

	if err := database.DB.Model(&existingTask).Updates(updates).Error; err != nil {
		return models.TaskResponse{}, err
	}

	// Reload updated task
	if err := database.DB.First(&existingTask, id).Error; err != nil {
		return models.TaskResponse{}, err
	}

	return models.TaskResponse{
		ID:          existingTask.ID,
		Title:       existingTask.Title,
		Description: existingTask.Description,
		Done:        existingTask.Done,
		UserID:      existingTask.UserID,
	}, nil
}

func DeleteTask(id uint) error {
	var task models.Task

	err := database.DB.First(&task, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTaskNotFound
	}

	if err != nil {
		return err
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		return err
	}

	return nil
}
