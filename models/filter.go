package models

type TaskFilter struct {
	Done   *bool
	UserID uint
	Sort   string
	Order  string
	Search string
}
