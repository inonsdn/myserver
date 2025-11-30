package main

import (
	"userms/internal/router"

	"github.com/inonsdn/myserver/http_con"
)

func main() {
	con := http_con.NewHandler()
	con.RegisterRoute(router.MainRoute{})

	go con.Run(":8081")

	con.WaitAndGetStatus()
}
