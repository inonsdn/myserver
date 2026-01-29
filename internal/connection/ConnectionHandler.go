package connection

import (
	"context"
	"fmt"
	"log/slog"
	"myserver/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	// get route from difinition in routeHandlePath
	pathToHandler := getRoutes()

	slog.Info(fmt.Sprintf("Register %d routes", len(pathToHandler)))

	// loop over path to set handle to http
	for path, handler := range pathToHandler {

		// set handle
		http.Handle(path, makeHandler(handler))
		slog.Debug(fmt.Sprintf("Found path for register %s", path))
	}
}

func runServer(server *http.Server, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	// run serve
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Stop service with error")
	}
}

func (c *ConnectionHandler) RunServe() {
	// get address to serve from config
	addr := c.config.GetAddr()
	slog.Info(fmt.Sprintf("Run serve address %s", addr))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// get context for cancel it when background is done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := &http.Server{
		Addr: addr,
	}

	serveDone := make(chan struct{})
	go runServer(server, serveDone)

	select {
	// cancel with
	case <-sigChan:
		server.Shutdown(ctx)
	case <-serveDone:
		slog.Info("Server is shutting down")
	}
}
