package main

import (
	"fmt"
	"log/slog"
	"myserver/internal/config"
	"myserver/internal/connection"
	"myserver/internal/database"
	"os"
)

func testDbCon() {
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

	userCon := dbHandler.GetUserConnection()
	userId := userCon.CreateNewUser("nonser")

	slog.Info(fmt.Sprintf("Create user and got id %s", userId))
	allUsers := userCon.GetAllUser()
	slog.Info(fmt.Sprintf("Get all user: %v", allUsers))
}

func main() {
	// setting logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(logger)

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
