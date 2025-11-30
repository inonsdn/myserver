package httpcon

import (
	"fmt"
	"net/http"
)

type HttpConnection struct {
	channel chan int
	path    []Path
}

func InitPathFromConfig(h *HttpConnection, paths []Path) {
	for _, path := range paths {
		http.HandleFunc(path.Name, path.Callback)
		h.path = append(h.path, path)
	}
}

// func InitPath(h *HttpConnection, pathHandler PathHandler) {
// 	http.HandleFunc(pathHandler.GetPathName(), pathHandler.Execute)
// 	h.path = append(h.path, pathHandler)
// }

func (h *HttpConnection) Run() error {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Got error while running, stop")
	}
	return err
}

func (h *HttpConnection) RunForever() {
	fmt.Printf("Running and waiting\n")
	err := h.Run()
	if err != nil {
		fmt.Println("Got error while running, stop")
		h.channel <- -1
	} else {
		h.channel <- 0
	}
}

func (h *HttpConnection) WaitStatus() int {
	status := <-h.channel
	fmt.Printf("Running got status: %d\n", status)
	return status
}
