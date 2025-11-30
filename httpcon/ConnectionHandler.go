package httpcon

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	route   *gin.Engine
	sigChan chan int
}

func InitHandler() (*ConnectionHandler, error) {
	route := gin.Default()
	var err error
	for _, httpPath := range ConfigPaths {
		if httpPath.method == "GET" {
			route.GET(httpPath.name, httpPath.callback)
		} else if httpPath.method == "POST" {
			route.POST(httpPath.name, httpPath.callback)
		} else {
			err = fmt.Errorf("")
		}
	}

	return &ConnectionHandler{
		route:   route,
		sigChan: make(chan int, 0),
	}, err
}

func (c *ConnectionHandler) WaitAndGetStatus() int {
	return <-c.sigChan
}

func (c *ConnectionHandler) Run() {
	err := c.route.Run()
	if err != nil {
		fmt.Println("Found error")
		c.sigChan <- -1
	}
	c.sigChan <- 0
}
