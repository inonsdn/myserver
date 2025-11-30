package main

import (
	"github.com/inonsdn/http_con"
	"github.com/inonsdn/myserver/gateway/internal/router"
)

func main() {
	con := http_con.NewHandler()
	con.RegisterRoute(router.MainRoute{})

	go con.Run(":8080")

	con.WaitAndGetStatus()
}
