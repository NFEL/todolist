package middlewares

import (
	"bytes"
	"graph-interview/internal/cfg"
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func AccessLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Cfg.Server.AccessLog {
			return
		}
		// Process request
		buf, _ := io.ReadAll(c.Request.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = rdr2

		logger.InfoContext(c.Request.Context(), "accesslog", "method", c.Request.Method, "route", c.FullPath(), "req", rdr1)

		writer := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		// Process request
		c.Next()

		defer func() {
			logger.InfoContext(c.Request.Context(), "accesslog", "method", c.Request.Method, "status", c.Writer.Status(), "route", c.FullPath(), "resp", writer.body.String()[:25])
		}()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
