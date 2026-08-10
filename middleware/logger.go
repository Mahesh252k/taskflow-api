package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Before request", c.Request.Method, c.Request.URL.Path)

		c.Next()

		fmt.Println("After request", c.Writer.Status())
	}
}
