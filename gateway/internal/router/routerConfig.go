package router

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	tokenTimestamp = int64(3600) // 1 hours for timestamp of token valid
)

type MainRoute struct{}

func (m MainRoute) RegisterRoute(r *gin.Engine) {
	RegisterGenericRoute(r)
	RegisterUserMSRoute(r)
}

func RegisterGenericRoute(r *gin.Engine) {
	r.GET("/ping", pong)
}

func RegisterUserMSRoute(r *gin.Engine) {
	// get proxy to user ms
	target, err := url.Parse("http://localhost:8081")
	if err != nil {
		fmt.Println("Got error from register")
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	// register this will forward to user ms
	r.POST("/login", forwardToUserMs(proxy))
	r.POST("/createUser", forwardToUserMs(proxy))

	apiGroup := r.Group("/api")
	apiGroup.Use(AuthorizeJWT())
	{
		apiGroup.GET("/getUserInfo", forwardToUserMs(proxy))
	}
}
