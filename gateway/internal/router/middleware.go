package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	BEARER_SCHEMA = "Bearer "
)

func getTokenFromHeader(authHeader string) string {
	// init token value
	token := ""

	// if there is no header
	if authHeader == "" {
		return token
	}

	// split header to 2 elements by space
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return token
	}

	// check prefix of authorize
	prefix := strings.ToLower(parts[0])
	if prefix != "bearer" {
		return token
	}
	token = parts[1]

	return token
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get header of authorization
		authHeader := c.GetHeader("Authorization")

		token := getTokenFromHeader(authHeader)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing token, reauthorize again",
			})
		}

		c.Request.Header.Set("User-Id", "1")
		c.Set("userId", "1")
		c.Next()

	}
}
