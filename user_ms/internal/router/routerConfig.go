package router

import (
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("super-secret-demo") // must match gateway

const (
	tokenTimestamp = int64(3600) // 1 hours for timestamp of token valid

)

type MainRoute struct{}

func (m MainRoute) RegisterRoute(r *gin.Engine) {
	RegisterGenericRoute(r)
}

func RegisterGenericRoute(r *gin.Engine) {
	r.GET("/ping", pong)
	r.POST("/login", login)
}

func RegisterAuthRoute(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use()
	{
		apiGroup.POST("/getUserInfo", login)
	}
}
