package connection

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIdKey contextKey = "userId"

type middleware func(RouteHandlerFunc) RouteHandlerFunc

// middleware for connection handler which will execute before handler function
// to check header and authenticate user with jwt token
func authMiddleware(secret []byte, next RouteHandlerFunc) RouteHandlerFunc {
	return func(rh *RouteHandler) error {

		// get header and check if it has "Authorization" header
		authHeader := rh.r.Header.Get("Authorization")
		if authHeader == "" {
			rh.ResponseError(http.StatusUnauthorized, "Missing Authorization header")
			return nil
		}

		// check if header is in format "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") || parts[1] == "" {
			rh.ResponseError(http.StatusUnauthorized, "Invalid Authorization header format")
			return nil
		}
		tokenStr := parts[1]

		// parse and verify JWT token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			rh.ResponseError(http.StatusUnauthorized, "Invalid token")
			return nil
		}

		// extract userId from claims and attach to request context
		userId, ok := claims["userId"].(string)
		if !ok || userId == "" {
			rh.ResponseError(http.StatusUnauthorized, "Invalid token claims")
			return nil
		}

		ctx := context.WithValue(rh.r.Context(), userIdKey, userId)
		rh.r = rh.r.WithContext(ctx)

		// if token is valid, execute next handler function
		return next(rh)
	}
}
