package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model

	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`

	UserID uint `json:"user_id"`
	User   User `json:"user"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type TaskListResponse struct {
	Data       []Task     `json:"data"`
	Pagination Pagination `json:"pagination"`
}
