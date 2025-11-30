package main

import (
	"github.com/inonsdn/myserver/httpcon"
)

func main() {
	con, err := httpcon.InitHandlerWithGroup()
	if err != nil {
		return
	}
	go con.Run()

	con.WaitAndGetStatus()
}
