package main

import (
	"github.com/inonsdn/myserver/gateway/internal/router"
	"github.com/inonsdn/myserver/http_con"
)

func main() {
	routerHandler := router.MainRoute{}
	con := http_con.NewHandler(routerHandler)
	con.RegisterRoute()

	go con.Run(":8080")

	con.WaitAndGetStatus()
}
