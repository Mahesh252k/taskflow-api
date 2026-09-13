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

const (
	requestLimit = 5
	window       = time.Minute
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

		if time.Since(limit.StartedAt) >= window {
			limit.Count = 0
			limit.StartedAt = time.Now()
		}

		limit.Count++

		if limit.Count > requestLimit {
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

func CleanupRateLimits() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()

		for userID, limit := range users {
			if time.Since(limit.StartedAt) >= window {
				delete(users, userID)
			}
		}

		mu.Unlock()
	}
}
