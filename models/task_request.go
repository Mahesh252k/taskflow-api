package models

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id" binding:"required"`
}
type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}
