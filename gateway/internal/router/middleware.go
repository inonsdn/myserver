package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token, reauthorize again",
			})
			return
		}

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			// TODO: replace with your actual signing key or key lookup
			return []byte("super-secret-demo"), nil
		})
		if err != nil || parsedToken == nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		expUnix, ok := claims["expiredUnix"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "not found expiration time",
			})
			return
		}

		expTime := time.Unix(int64(expUnix), 0)
		loginTime := time.Now()
		if loginTime.Sub(expTime) > time.Duration(tokenTimestamp) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok {
			userId = "-1"
		}
		c.Request.Header.Set("userId", userId)
		c.Set("userId", userId)
		c.Next()
	}
}
