package main

import (
	"log/slog"
	"myserver/internal/config"
	"myserver/internal/connection"
	"myserver/internal/database"
)

func main() {

	// load config
	serverConfig := config.LoadConfig()
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		slog.Error("Got error when load config")
		slog.Error(err.Error())
		return
	}
	pgExecutor := database.NewPGExecutor(dbConfig)
	dbHandler := database.Connect(pgExecutor)

	if dbHandler == nil {
		return
	}

	// construct http connection handler
	handler := connection.NewConnectionHandler(serverConfig, dbHandler)

	// register route handler before run
	handler.RegisterRoute()

	// run forever loop
	handler.RunServe()

	// error handling
}
