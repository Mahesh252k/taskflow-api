package handlers

import (
	"errors"
	"net/http"

	"taskflow-api/services"

	"github.com/gin-gonic/gin"
)

func HandleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrTitleRequired),
		errors.Is(err, services.ErrUserRequired):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, services.ErrTaskNotFound),
		errors.Is(err, services.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
