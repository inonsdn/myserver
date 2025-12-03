package main

import (
	dbcon "userms/internal/dbCon"
	"userms/internal/router"

	"github.com/inonsdn/myserver/http_con"
)

func main() {
	// TODO: setting this
	localConfig := dbcon.DbConfig{
		Host:     "localhost",
		Port:     8000,
		User:     "admin",
		Password: "test",
		DBName:   "userdb",
	}
	userRouterHandler := router.NewRouterHandler(&localConfig)
	con := http_con.NewHandler(userRouterHandler)
	con.RegisterRoute()

	go con.Run(":8081")

	con.WaitAndGetStatus()
}
