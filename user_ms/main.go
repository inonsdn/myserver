package main

import (
	"userms/internal/router"

	"github.com/inonsdn/myserver/http_con"
)

func main() {
	// define config function to override default config
	configFuncs := []router.OptsFunc{
		tokenPeriodTimestamp(int64(120)),
	}

	userRouterHandler := router.NewRouterHandler(configFuncs...)
	con := http_con.NewHandler(userRouterHandler)
	con.RegisterRoute()

	go con.Run(":8081")

	con.WaitAndGetStatus()
}
