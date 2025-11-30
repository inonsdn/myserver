package router

import (
	"github.com/gin-gonic/gin"
)

type MainRoute struct{}

func (m MainRoute) RegisterRoute(r *gin.Engine) {
	RegisterGenericRoute(r)
	RegisterAuthRoute(r)
}

func RegisterGenericRoute(r *gin.Engine) {
	r.GET("/ping", pong)
}

func RegisterAuthRoute(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use(AuthorizeJWT())
	{
		apiGroup.GET("/getUserInfo")
	}
}
