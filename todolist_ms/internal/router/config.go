package router

import "github.com/gin-gonic/gin"

type RouteRegisterFunc func(r *gin.Engine)

type RouteHandler struct {
	route *gin.Engine
}

func Must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}

func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

func (r *RouteHandler) RegisterRoute(f ...RouteRegisterFunc) {
	for _, fCall := range f {
		fCall(r.route)
	}
}
