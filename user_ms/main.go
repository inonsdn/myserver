package main

import (
	dbcon "userms/internal/dbCon"
	"userms/internal/router"

	"github.com/inonsdn/myserver/http_con"
)

func main() {
	// TODO: setting this
	localConfig := dbcon.DbConfig{
		// Host     string
		// Port     int
		// User     string
		// Password string
		// DBName   string
	}
	userRouterHandler := router.NewRouterHandler(&localConfig)
	con := http_con.NewHandler(userRouterHandler)
	con.RegisterRoute()

	go con.Run(":8081")

	con.WaitAndGetStatus()
}
