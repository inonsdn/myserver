package connection

import (
	"fmt"
	"log/slog"
	"myserver/internal/config"
	"net/http"
)

type ConnectionHandler struct {
	config *config.Config
}

func NewConnectionHandler(config *config.Config) *ConnectionHandler {
	slog.Info("Create connection")
	return &ConnectionHandler{
		config: config,
	}
}

func (c *ConnectionHandler) RegisterRoute() {
	pathToHandler := getRoutes()
	slog.Info(fmt.Sprintf("Register %d routes", len(pathToHandler)))
	for path, handler := range pathToHandler {
		http.Handle(path, makeHandler(handler))
		slog.Debug(fmt.Sprintf("Found path for register %s", path))
	}
}

func (c *ConnectionHandler) RunServe() {

	addr := c.config.GetAddr()
	slog.Info(fmt.Sprintf("Run serve address %s", addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("Stop service with error")
	}
}
