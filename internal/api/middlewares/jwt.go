package middlewares

import (
	"context"
	"errors"
	"graph-interview/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func bearerFromHeader(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

func AuthMiddleware(authSrv *services.AuthService, r *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, _ := c.Cookie("access_token")
		if tokenStr == "" {
			tokenStr = bearerFromHeader(c)
		}
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := authSrv.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx := context.Background()
		// return r.Client.Set(ctx, key, userID, time.Until(exp)).Err()
		if cmd := r.Get(ctx, "access:"+claims.ID); cmd.Err() != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}

func MustCookie(c *gin.Context, name string) (string, error) {
	val, err := c.Cookie(name)
	if err != nil || val == "" {
		return "", errors.New("missing cookie: " + name)
	}
	return val, nil
}
