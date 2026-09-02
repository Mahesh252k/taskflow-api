package handlers

import "github.com/gin-gonic/gin"

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
