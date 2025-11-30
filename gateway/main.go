package main

import (
	"github.com/inonsdn/http_con"
	"github.com/inonsdn/myserver/gateway/internal/config"
)

func main() {
	con := http_con.NewHandler()
	con.RegisterRoute(config.MainRoute{})

	go con.Run(":8080")

	con.WaitAndGetStatus()
}
