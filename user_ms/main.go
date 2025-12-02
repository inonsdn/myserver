package main

import (
	"userms/internal/router"

	"github.com/inonsdn/myserver/http_con"
)

func main() {
	userRouterHandler := router.NewRouterHandler()
	con := http_con.NewHandler(userRouterHandler)
	con.RegisterRoute()

	go con.Run(":8081")

	con.WaitAndGetStatus()
}
