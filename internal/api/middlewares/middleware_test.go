package middlewares

import (
	"graph-interview/internal/cfg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestJSONMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(JSONMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestJSONMiddleware_StaticPath(t *testing.T) {
	r := gin.New()
	r.Use(JSONMiddleware())
	r.GET("/public/static/test.js", func(c *gin.Context) {
		c.String(200, "content")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/public/static/test.js", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Static path should NOT have json content type forced
	assert.NotContains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestCorsMiddleware(t *testing.T) {
	corsCfg := cfg.CorsCfg{
		Origins:        []string{"http://localhost:3000"},
		Methods:        []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	r := gin.New()
	r.Use(CorsMiddleware(corsCfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCorsMiddleware_Options(t *testing.T) {
	corsCfg := cfg.CorsCfg{
		Origins:        []string{"http://localhost:3000"},
		Methods:        []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	}

	r := gin.New()
	r.Use(CorsMiddleware(corsCfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestPrometheusMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(PrometheusMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBearerFromHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"valid bearer", "Bearer abc123", "abc123"},
		{"no bearer prefix", "abc123", ""},
		{"empty header", "", ""},
		{"bearer with space", "Bearer token with space", "token with space"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)
			if tt.header != "" {
				c.Request.Header.Set("Authorization", tt.header)
			}
			result := bearerFromHeader(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMustCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// No cookie set
	val, err := MustCookie(c, "test_cookie")
	assert.Error(t, err)
	assert.Empty(t, val)
}

func TestMustCookie_WithCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "test_cookie", Value: "test_value"})

	val, err := MustCookie(c, "test_cookie")
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)
}
