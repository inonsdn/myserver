package main

import (
	"log/slog"
	"myserver/internal/config"
	"myserver/internal/connection"
	"os"
)

func main() {
	// setting logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(logger)

	// load config
	config := config.LoadConfig()

	// construct http connection handler
	handler := connection.NewConnectionHandler(config)

	// register route handler before run
	handler.RegisterRoute()

	// run forever loop
	handler.RunServe()

	// error handling
}
