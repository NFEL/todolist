package middlewares

import (
	"fmt"
	"graph-interview/internal/cfg"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
)

func ProfilingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Cfg != nil && cfg.Cfg.Verbose {
			startAt := time.Now().UTC()
			c.Next()
			finishedAt := time.Now().UTC()
			duration := finishedAt.Sub(startAt)

			switch true {
			case duration.Seconds() > 0.3:
				fmt.Println(colorRed, "===EXECUTION TIME: ", duration, colorReset)
			case duration.Seconds() > 1:
				fmt.Println(colorYellow, "===EXECUTION TIME: ", duration, colorReset)
			default:
				fmt.Println(colorCyan, "===EXECUTION TIME: ", duration, colorReset)
			}
		} else {
			startAt := time.Now().UTC()
			c.Next()
			finishedAt := time.Now().UTC()
			duration := finishedAt.Sub(startAt)
			logger.Info("execution time", "duration", duration)
		}
	}
}
