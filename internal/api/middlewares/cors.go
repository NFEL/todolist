package middlewares

import (
	"graph-interview/internal/cfg"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(cfg cfg.CorsCfg) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(cfg.Origins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.Methods, ","))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
