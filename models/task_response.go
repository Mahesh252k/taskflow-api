package models

type TaskResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	UserID      uint   `json:"user_id"`
}
