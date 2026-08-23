package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"taskflow-api/models"
	"taskflow-api/services"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var req models.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := services.RegisterUser(req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})

}

func LoginUser(c *gin.Context) {
	var req models.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := services.LoginUser(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged in successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetTasks(c *gin.Context) {
	pageString := c.DefaultQuery("page", "1")
	limitString := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageString)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "page must be a positive integer",
		})
		return
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be a positive integer between 1 and 100",
		})
		return
	}

	sort := strings.ToLower(c.DefaultQuery("sort", "id"))

	allowedSorts := map[string]bool{
		"id":         true,
		"title":      true,
		"created_at": true,
	}

	if !allowedSorts[sort] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid sort field",
		})
		return
	}

	order := strings.ToLower(c.DefaultQuery("order", "asc"))

	allowedOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedOrders[order] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order value",
		})
		return
	}

	search := strings.TrimSpace(c.Query("search"))

	filter := models.TaskFilter{
		Sort:   sort,
		Order:  order,
		Search: search,
	}

	if doneString, exists := c.GetQuery("done"); exists {
		done, err := strconv.ParseBool(doneString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "done must be true or false",
			})
			return
		}

		filter.Done = &done
	}

	if userIDString, exists := c.GetQuery("user_id"); exists {
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil || userID < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user_id must be a positive integer",
			})
			return
		}

		filter.UserID = uint(userID)
	}

	response, err := services.GetTasks(page, limit, filter)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func CreateTask(c *gin.Context) {
	var createTaskRequest models.CreateTaskRequest

	if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation failed: " + err.Error(),
		})
		return
	}

	createdTask, err := services.CreateTask(models.Task{
		Title:       createTaskRequest.Title,
		Description: createTaskRequest.Description,
		UserID:      createTaskRequest.UserID,
	})
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := models.TaskResponse{
		ID:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Done:        createdTask.Done,
		UserID:      createdTask.UserID,
	}

	c.JSON(http.StatusCreated, response)
}

func GetTaskByID(c *gin.Context) {
	id, err := parseTaskID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	task, err := services.GetTaskByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id, err := parseTaskID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	var updateTaskRequest models.UpdateTaskRequest

	if err := c.ShouldBindJSON(&updateTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation failed: " + err.Error(),
		})
		return
	}

	if updateTaskRequest.Title == nil &&
		updateTaskRequest.Description == nil &&
		updateTaskRequest.Done == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "at least one field (title, description, done) must be provided for update",
		})
		return
	}

	updatedTask, err := services.UpdateTask(id, updateTaskRequest)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func DeleteTask(c *gin.Context) {
	id, err := parseTaskID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	err = services.DeleteTask(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}

func parseTaskID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid task ID")
	}

	return uint(id), nil
}
