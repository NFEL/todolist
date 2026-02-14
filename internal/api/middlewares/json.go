package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.URL.Path, "/public/static/") {
			c.Writer.Header().Set("Content-Type", "application/json")
		}
		c.Next()
	}
}
