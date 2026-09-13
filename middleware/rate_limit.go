package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type userLimit struct {
	Count     int
	StartedAt time.Time
}

var (
	users = make(map[uint]*userLimit)
	mu    sync.Mutex
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		mu.Lock()

		limit, exists := users[userID]

		if !exists {
			limit = &userLimit{
				Count:     0,
				StartedAt: time.Now(),
			}
			users[userID] = limit
		}

		if time.Since(limit.StartedAt) >= time.Minute {
			limit.Count = 0
			limit.StartedAt = time.Now()
		}

		limit.Count++

		if limit.Count > 5 {
			mu.Unlock()

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}

		mu.Unlock()

		c.Next()
	}
}
