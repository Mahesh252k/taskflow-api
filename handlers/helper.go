package handlers

import (
	"errors"
	"strconv"
	"strings"
	"taskflow-api/models"

	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}

func parseTaskQuery(c *gin.Context) (int, int, models.TaskFilter, error) {
	pageString := c.DefaultQuery("page", "1")
	limitString := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageString)
	if err != nil || page < 1 {
		return 0, 0, models.TaskFilter{}, errors.New("page must be a positive integer")
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil || limit < 1 || limit > 100 {
		return 0, 0, models.TaskFilter{}, errors.New("limit must be a positive integer between 1 and 100")
	}

	sort := strings.ToLower(c.DefaultQuery("sort", "id"))

	allowedSorts := map[string]bool{
		"id":         true,
		"title":      true,
		"created_at": true,
	}

	if !allowedSorts[sort] {
		return 0, 0, models.TaskFilter{}, errors.New("invalid sort field")
	}

	order := strings.ToLower(c.DefaultQuery("order", "asc"))

	allowedOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedOrders[order] {
		return 0, 0, models.TaskFilter{}, errors.New("invalid order value")
	}

	filter := models.TaskFilter{
		Sort:   sort,
		Order:  order,
		Search: strings.TrimSpace(c.Query("search")),
	}

	if doneString, exists := c.GetQuery("done"); exists {
		done, err := strconv.ParseBool(doneString)
		if err != nil {
			return 0, 0, models.TaskFilter{}, errors.New("done must be true or false")
		}

		filter.Done = &done
	}

	return page, limit, filter, nil
}
