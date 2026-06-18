package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware(keys map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			key = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		}

		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			return
		}

		app := c.GetHeader("X-App")
		if app == "" {
			app = c.Query("app")
		}

		expected, ok := keys[strings.ToLower(app)]
		if !ok || expected != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}

		c.Set("app", strings.ToLower(app))
		c.Next()
	}
}