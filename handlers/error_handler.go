package handlers

import (
	"errors"
	"log"
	"net/http"

	"taskflow-api/services"

	"github.com/gin-gonic/gin"
)

func handleServiceError(c *gin.Context, err error) {
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
		log.Print("unexpected service error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
