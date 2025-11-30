package main

import (
	"userms/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/inonsdn/http_con"
)

var httpConfigs = []http_con.HttpGroupPath{
	{
		Name: "",
		Paths: []http_con.HttpPath{
			{
				Name:     "/getUserInfo",
				Callback: getUserInfo,
				Method:   http_con.RouteMethod_GET,
			},
		},
	},
}

func getUserInfo(c *gin.Context) {

}

func main() {
	con := http_con.NewHandler()
	con.RegisterRoute(router.MainRoute{})

	go con.Run(":8081")

	con.WaitAndGetStatus()
}
