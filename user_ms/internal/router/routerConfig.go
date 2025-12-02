package router

import (
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("super-secret-demo") // must match gateway

const (
	tokenTimestamp = int64(3600) // 1 hours for timestamp of token valid

)

func (u *UserRouterHandler) GetStruct() *UserRouterHandler {
	return u
}

func (u *UserRouterHandler) RegisterRoute(r *gin.Engine) {
	u.registerGenericRoute(r)
	u.registerAuthRoute(r)
}

func (u *UserRouterHandler) registerGenericRoute(r *gin.Engine) {
	r.GET("/ping", pong)
	r.POST("/login", u.login)
}

func (u *UserRouterHandler) registerAuthRoute(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use()
	{
		apiGroup.GET("/getUserInfo", u.getUserInfo)
	}
}
