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

	// Count AFTER filters, BEFORE pagination
	if err := query.Count(&total).Error; err != nil {
		return models.TaskListResponse{}, err
	}

	allowedSorts := map[string]bool{
		"id":         true,
		"title":      true,
		"created_at": true,
	}

	allowedOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedSorts[filter.Sort] {
		filter.Sort = "id"
	}

	if !allowedOrders[filter.Order] {
		filter.Order = "asc"
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

func GetTaskByID(id int) (models.Task, error) {
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

func UpdateTask(id int, updated models.Task) (models.Task, error) {
	var existingTask models.Task

	err := database.DB.First(&existingTask, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, ErrTaskNotFound
	}

	if err != nil {
		return models.Task{}, err
	}

	updated.Title = strings.TrimSpace(updated.Title)
	updated.Description = strings.TrimSpace(updated.Description)

	if updated.Title == "" {
		return models.Task{}, ErrTitleRequired
	}

	existingTask.Title = updated.Title
	existingTask.Description = updated.Description
	existingTask.Done = updated.Done

	if updated.UserID != 0 {
		var user models.User

		err := database.DB.First(&user, updated.UserID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Task{}, ErrUserNotFound
		}

		if err != nil {
			return models.Task{}, err
		}

		existingTask.UserID = updated.UserID
	}

	if err := database.DB.Save(&existingTask).Error; err != nil {
		return models.Task{}, err
	}

	if err := database.DB.
		Preload("User").
		First(&existingTask, existingTask.ID).Error; err != nil {
		return models.Task{}, err
	}

	return existingTask, nil
}

func DeleteTask(id int) error {
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
