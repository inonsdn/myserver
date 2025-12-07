package router

import (
	dbcon "userms/internal/dbCon"

	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("super-secret-demo") // must match gateway

const (
	tokenTimestamp = int64(3600) // 1 hours for timestamp of token valid

)

type Options struct {
	dbcon.DbConfig
	TokenPeriodTimestamp int64
	JwtSecret            []byte
}

type OptsFunc func(*Options)

func defaultOptions() Options {
	return Options{
		DbConfig: dbcon.DbConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "rootPass",
			DBName:   "userdb",
		},
		TokenPeriodTimestamp: int64(3600),
		JwtSecret:            []byte("super-secret-demo"),
	}
}

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
	r.POST("/createUser", u.createUser)
}

func (u *UserRouterHandler) registerAuthRoute(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use()
	{
		apiGroup.GET("/getUserInfo", u.getUserInfo)
	}
}
